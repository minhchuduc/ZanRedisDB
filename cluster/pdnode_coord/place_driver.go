package pdnode_coord

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/spaolacci/murmur3"
	"github.com/youzan/ZanRedisDB/cluster"
	"github.com/youzan/ZanRedisDB/common"
)

var (
	ErrBalanceNodeUnavailable = errors.New("can not find a node to be balanced")
	ErrClusterBalanceRunning  = errors.New("another balance is running, should wait")
)

type balanceOpLevel int

func splitNamespacePartitionID(namespaceFullName string) (string, int, error) {
	partIndex := strings.LastIndex(namespaceFullName, "-")
	if partIndex == -1 {
		return "", 0, fmt.Errorf("invalid namespace full name: %v", namespaceFullName)
	}
	namespaceName := namespaceFullName[:partIndex]
	partitionID, err := strconv.Atoi(namespaceFullName[partIndex+1:])
	return namespaceName, partitionID, err
}

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getNodeNameList(currentNodes map[string]cluster.NodeInfo) []SortableStrings {
	nodeNameMap := make(map[string]SortableStrings)
	dcInfoList := make(SortableStrings, 0)
	for nid, ninfo := range currentNodes {
		dcInfo := ""
		dc, ok := ninfo.Tags[cluster.DCInfoTag]
		if ok {
			dcInfo, _ = dc.(string)
		}
		nodeNameMap[dcInfo] = append(nodeNameMap[dcInfo], nid)
	}
	for dcInfo := range nodeNameMap {
		dcInfoList = append(dcInfoList, dcInfo)
	}
	sort.Sort(dcInfoList)
	nodeNameList := make([]SortableStrings, 0, len(nodeNameMap))
	for _, dc := range dcInfoList {
		sort.Sort(nodeNameMap[dc])
		nodeNameList = append(nodeNameList, nodeNameMap[dc])
	}
	return nodeNameList
}

// query the raft peers if the nid already in the raft group for the namespace in all raft peers
func IsRaftNodeJoined(nsInfo *cluster.PartitionMetaInfo, nid string) (bool, error) {
	if len(nsInfo.RaftNodes) == 0 {
		return false, nil
	}
	var lastErr error
	for _, remoteNode := range nsInfo.GetISR() {
		if remoteNode == nid {
			continue
		}
		nip, _, _, httpPort := cluster.ExtractNodeInfoFromID(remoteNode)
		var rsp []*common.MemberInfo
		_, err := common.APIRequest("GET",
			"http://"+net.JoinHostPort(nip, httpPort)+common.APIGetMembers+"/"+nsInfo.GetDesp(),
			nil, time.Second*3, &rsp)
		if err != nil {
			cluster.CoordLog().Infof("failed (%v) to get members for namespace %v: %v", nip, nsInfo.GetDesp(), err)
			lastErr = err
			continue
		}

		for _, m := range rsp {
			if m.NodeID == cluster.ExtractRegIDFromGenID(nid) && m.ID == nsInfo.RaftIDs[nid] {
				cluster.CoordLog().Infof("namespace %v node %v is still in raft from %v ", nsInfo.GetDesp(), nid, remoteNode)
				return true, nil
			}
		}
	}
	return false, lastErr
}

// query the raft peers if all the isr nodes already in the raft group for the namespace and all logs synced in all peers
func IsAllISRFullReady(nsInfo *cluster.PartitionMetaInfo) (bool, error) {
	for _, nid := range nsInfo.GetISR() {
		ok, err := IsRaftNodeFullReady(nsInfo, nid)
		if err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}

// query the raft peers if the nid already in the raft group for the namespace and all logs synced in all peers
func IsRaftNodeFullReady(nsInfo *cluster.PartitionMetaInfo, nid string) (bool, error) {
	if len(nsInfo.RaftNodes) == 0 {
		return false, nil
	}
	for _, remoteNode := range nsInfo.GetISR() {
		nip, _, _, httpPort := cluster.ExtractNodeInfoFromID(remoteNode)
		var rsp []*common.MemberInfo
		_, err := common.APIRequest("GET",
			"http://"+net.JoinHostPort(nip, httpPort)+common.APIGetMembers+"/"+nsInfo.GetDesp(),
			nil, time.Second*3, &rsp)
		if err != nil {
			cluster.CoordLog().Infof("failed (%v) to get members for namespace %v: %v", nip, nsInfo.GetDesp(), err)
			return false, err
		}

		found := false
		for _, m := range rsp {
			if m.NodeID == cluster.ExtractRegIDFromGenID(nid) && m.ID == nsInfo.RaftIDs[nid] {
				found = true
				break
			}
		}
		if !found {
			cluster.CoordLog().Infof("raft %v not found in the node (%v) members for namespace %v", nid, nip, nsInfo.GetDesp())
			return false, nil
		}
		_, err = common.APIRequest("GET",
			"http://"+net.JoinHostPort(nip, httpPort)+common.APIIsRaftSynced+"/"+nsInfo.GetDesp(),
			nil, time.Second*5, nil)
		if err != nil {
			cluster.CoordLog().Infof("failed (%v) to check sync state for namespace %v: %v", nip, nsInfo.GetDesp(), err)
			return false, err
		}
	}
	return true, nil
}

// check if all nodes in list is raft synced for the namespace-partition
func IsRaftNodeSynced(nsInfo *cluster.PartitionMetaInfo, nid string) (bool, error) {
	nip, _, _, httpPort := cluster.ExtractNodeInfoFromID(nid)
	_, err := common.APIRequest("GET",
		"http://"+net.JoinHostPort(nip, httpPort)+common.APIIsRaftSynced+"/"+nsInfo.GetDesp(),
		nil, time.Second*10, nil)
	if err != nil {
		cluster.CoordLog().Infof("failed (%v) to check sync state for namespace %v: %v", nip, nsInfo.GetDesp(), err)
		return false, err
	}
	return true, nil
}

type DataPlacement struct {
	balanceInterval [2]int32
	balanceVer      string
	pdCoord         *PDCoordinator
}

func NewDataPlacement(coord *PDCoordinator) *DataPlacement {
	return &DataPlacement{
		pdCoord:         coord,
		balanceInterval: [2]int32{2, 4},
	}
}

func (dp *DataPlacement) SetBalanceInterval(start int, end int) {
	if start == end && start == 0 {
		return
	}
	atomic.StoreInt32(&dp.balanceInterval[0], int32(start))
	atomic.StoreInt32(&dp.balanceInterval[1], int32(end))
}

func (dp *DataPlacement) DoBalance(monitorChan chan struct{}) {
	//check period for the data balance.
	ticker := time.NewTicker(balanceCheckInterval)
	defer func() {
		ticker.Stop()
		cluster.CoordLog().Infof("balance check exit.")
	}()
	for {
		select {
		case <-monitorChan:
			return
		case <-ticker.C:
			// only balance at given interval
			if time.Now().Hour() > int(atomic.LoadInt32(&dp.balanceInterval[1])) ||
				time.Now().Hour() < int(atomic.LoadInt32(&dp.balanceInterval[0])) {
				continue
			}
			if !dp.pdCoord.IsMineLeader() {
				cluster.CoordLog().Infof("not leader while checking balance")
				continue
			}
			if !dp.pdCoord.IsClusterStable() {
				cluster.CoordLog().Infof("no balance since cluster is not stable while checking balance")
				continue
			}
			if !dp.pdCoord.AutoBalanceEnabled() {
				continue
			}
			cluster.CoordLog().Infof("begin checking balance of namespace data...")
			currentNodes := dp.pdCoord.getCurrentNodes(nil)
			validNum := len(currentNodes)
			if validNum < 2 {
				continue
			}
			dp.rebalanceNamespace(monitorChan)
		}
	}
}

func (dp *DataPlacement) addNodeToNamespaceAndWaitReady(monitorChan chan struct{}, namespaceInfo *cluster.PartitionMetaInfo,
	nodeNameList []SortableStrings) (*cluster.PartitionMetaInfo, error) {
	retry := 0
	currentSelect := 0
	namespaceName := namespaceInfo.Name
	partitionID := namespaceInfo.Partition
	partitionNodes, coordErr := getRebalancedPartitionsFromNameList(
		namespaceInfo.Name,
		namespaceInfo.PartitionNum,
		namespaceInfo.Replica, nodeNameList, dp.balanceVer)
	if coordErr != nil {
		return namespaceInfo, coordErr.ToErrorType()
	}
	fullName := namespaceInfo.GetDesp()
	selectedCatchup := make([]string, 0)
	// we need add new catchup, choose from the the diff between the new isr and old isr
	for _, nid := range partitionNodes[namespaceInfo.Partition] {
		if cluster.FindSlice(namespaceInfo.RaftNodes, nid) != -1 {
			// already isr, ignore add catchup
			continue
		}
		selectedCatchup = append(selectedCatchup, nid)
	}
	var nInfo *cluster.PartitionMetaInfo
	var err error
	for {
		if currentSelect >= len(selectedCatchup) {
			cluster.CoordLog().Infof("currently no any node %v can be balanced for namespace: %v, expect isr: %v, nodes:%v",
				selectedCatchup, fullName, partitionNodes[partitionID], nodeNameList)
			return nInfo, ErrBalanceNodeUnavailable
		}
		nid := selectedCatchup[currentSelect]
		nInfo, err = dp.pdCoord.register.GetNamespacePartInfo(namespaceName, partitionID)
		if err != nil {
			cluster.CoordLog().Infof("failed to get namespace %v info: %v", fullName, err)
		} else {
			if inRaft, _ := IsRaftNodeFullReady(nInfo, nid); inRaft {
				break
			} else if cluster.FindSlice(nInfo.RaftNodes, nid) != -1 {
				// wait ready
				select {
				case <-monitorChan:
					return nInfo, errors.New("quiting")
				case <-time.After(time.Second * 5):
					cluster.CoordLog().Infof("node: %v is added for namespace %v , still waiting catchup", nid, nInfo.GetDesp())
				}
				continue
			} else {
				if ok, err := IsAllISRFullReady(nInfo); err != nil || !ok {
					cluster.CoordLog().Infof("namespace %v isr %v are not full ready while adding node", nInfo.GetDesp(), nInfo.RaftNodes)
					return nInfo, fmt.Errorf("namespace %v isr are not full ready", nInfo.GetDesp())
				}
				cluster.CoordLog().Infof("node: %v is added for namespace %v: (%v)", nid, nInfo.GetDesp(), nInfo.RaftNodes)
				coordErr = dp.pdCoord.addNamespaceToNode(nInfo, nid)
				if coordErr != nil {
					cluster.CoordLog().Infof("node: %v added for namespace %v (%v) failed: %v", nid,
						nInfo.GetDesp(), nInfo.RaftNodes, coordErr)
					currentSelect++
				}
			}
		}
		select {
		case <-monitorChan:
			return nInfo, errors.New("quiting")
		case <-time.After(time.Second * 5):
		}
		if retry > 5 {
			cluster.CoordLog().Infof("add catchup and wait timeout : %v", fullName)
			return nInfo, errors.New("wait timeout")
		}
		retry++
	}
	return nInfo, nil
}

func (dp *DataPlacement) getExcludeNodesForNamespace(namespaceInfo *cluster.PartitionMetaInfo) (map[string]struct{}, error) {
	excludeNodes := make(map[string]struct{})
	for _, v := range namespaceInfo.RaftNodes {
		excludeNodes[v] = struct{}{}
	}
	return excludeNodes, nil
}

func (dp *DataPlacement) allocNodeForNamespace(namespaceInfo *cluster.PartitionMetaInfo,
	currentNodes map[string]cluster.NodeInfo) (*cluster.NodeInfo, *cluster.CoordErr) {
	var chosenNode cluster.NodeInfo
	excludeNodes, commonErr := dp.getExcludeNodesForNamespace(namespaceInfo)
	if commonErr != nil {
		return nil, cluster.ErrRegisterServiceUnstable
	}

	partitionNodes, err := getRebalancedNamespacePartitions(
		namespaceInfo.Name,
		namespaceInfo.PartitionNum,
		namespaceInfo.Replica, currentNodes, dp.balanceVer)
	if err != nil {
		return nil, err
	}
	for _, nodeID := range partitionNodes[namespaceInfo.Partition] {
		if _, ok := excludeNodes[nodeID]; ok {
			continue
		}
		chosenNode = currentNodes[nodeID]
		break
	}
	if chosenNode.ID == "" {
		cluster.CoordLog().Infof("no more available node for namespace: %v, excluding nodes: %v, all nodes: %v",
			namespaceInfo.GetDesp(), excludeNodes, currentNodes)
		return nil, ErrNodeUnavailable
	}
	cluster.CoordLog().Infof("node %v is alloc for namespace: %v", chosenNode, namespaceInfo.GetDesp())
	return &chosenNode, nil
}

func (dp *DataPlacement) checkNamespaceNodeConflict(namespaceInfo *cluster.PartitionMetaInfo) bool {
	existSlaves := make(map[string]struct{})
	// isr should be different
	for _, id := range namespaceInfo.RaftNodes {
		if _, ok := existSlaves[id]; ok {
			return false
		}
		existSlaves[id] = struct{}{}
	}
	return true
}

func (dp *DataPlacement) allocNamespaceRaftNodes(ns string, currentNodes map[string]cluster.NodeInfo,
	replica int, partitionNum int, existPart map[int]*cluster.PartitionMetaInfo) ([]cluster.PartitionReplicaInfo, *cluster.CoordErr) {
	replicaList := make([]cluster.PartitionReplicaInfo, partitionNum)
	partitionNodes, err := getRebalancedNamespacePartitions(
		ns,
		partitionNum,
		replica, currentNodes, dp.balanceVer)
	if err != nil {
		return nil, err
	}
	for p := 0; p < partitionNum; p++ {
		var replicaInfo cluster.PartitionReplicaInfo
		if elem, ok := existPart[p]; ok {
			replicaInfo = elem.PartitionReplicaInfo
		} else {
			replicaInfo.RaftNodes = partitionNodes[p]
			replicaInfo.RaftIDs = make(map[string]uint64)
			replicaInfo.Removings = make(map[string]cluster.RemovingInfo)
			for _, nid := range replicaInfo.RaftNodes {
				replicaInfo.MaxRaftID++
				replicaInfo.RaftIDs[nid] = uint64(replicaInfo.MaxRaftID)
			}
		}
		replicaList[p] = replicaInfo
	}

	cluster.CoordLog().Infof("selected namespace replica list : %v", replicaList)
	return replicaList, nil
}

func (dp *DataPlacement) rebalanceNamespace(monitorChan chan struct{}) (bool, bool) {
	moved := false
	isAllBalanced := false
	if !atomic.CompareAndSwapInt32(&dp.pdCoord.balanceWaiting, 0, 1) {
		cluster.CoordLog().Infof("another balance is running, should wait")
		return moved, isAllBalanced
	}
	defer atomic.StoreInt32(&dp.pdCoord.balanceWaiting, 0)

	allNamespaces, _, err := dp.pdCoord.register.GetAllNamespaces()
	if err != nil {
		cluster.CoordLog().Infof("scan namespaces error: %v", err)
		return moved, isAllBalanced
	}
	namespaceList := make([]cluster.PartitionMetaInfo, 0)
	for _, parts := range allNamespaces {
		for _, p := range parts {
			namespaceList = append(namespaceList, *(p.GetCopy()))
		}
	}
	movedNamespace := ""
	isAllBalanced = true
	for _, namespaceInfo := range namespaceList {
		select {
		case <-monitorChan:
			return moved, false
		default:
		}
		if !dp.pdCoord.IsClusterStable() {
			return moved, false
		}
		if !dp.pdCoord.IsMineLeader() {
			return moved, true
		}
		if dp.pdCoord.hasRemovingNode() {
			return moved, false
		}
		// balance only one namespace once
		if movedNamespace != "" && movedNamespace != namespaceInfo.Name {
			continue
		}
		if len(namespaceInfo.Removings) > 0 {
			continue
		}
		if ok, err := IsAllISRFullReady(&namespaceInfo); err != nil || !ok {
			cluster.CoordLog().Infof("namespace %v isr is not full ready while balancing", namespaceInfo.GetDesp())
			continue
		}
		currentNodes := dp.pdCoord.getCurrentNodes(namespaceInfo.Tags)
		nodeNameList := getNodeNameList(currentNodes)
		cluster.CoordLog().Debugf("node name list: %v", nodeNameList)

		partitionNodes, err := getRebalancedNamespacePartitions(
			namespaceInfo.Name,
			namespaceInfo.PartitionNum,
			namespaceInfo.Replica, currentNodes, dp.balanceVer)
		if err != nil {
			isAllBalanced = false
			continue
		}
		cluster.CoordLog().Debugf("expected replicas : %v", partitionNodes)
		moveNodes := make([]string, 0)
		for _, nid := range namespaceInfo.GetISR() {
			found := false
			for _, expectedNode := range partitionNodes[namespaceInfo.Partition] {
				if nid == expectedNode {
					found = true
					break
				}
			}
			if !found {
				moveNodes = append(moveNodes, nid)
			}
		}
		for _, nid := range moveNodes {
			movedNamespace = namespaceInfo.Name
			cluster.CoordLog().Infof("node %v need move for namespace %v since %v not in expected isr list: %v", nid,
				namespaceInfo.GetDesp(), namespaceInfo.RaftNodes, partitionNodes[namespaceInfo.Partition])
			var err error
			var newInfo *cluster.PartitionMetaInfo
			if len(namespaceInfo.GetISR()) <= namespaceInfo.Replica {
				newInfo, err = dp.addNodeToNamespaceAndWaitReady(monitorChan, &namespaceInfo,
					nodeNameList)
			}
			if err != nil {
				return moved, false
			}
			if newInfo != nil {
				namespaceInfo = *newInfo
			}
			coordErr := dp.pdCoord.removeNamespaceFromNode(&namespaceInfo, nid)
			moved = true
			if coordErr != nil {
				cluster.CoordLog().Infof("namespace %v removing node: %v failed while balance",
					namespaceInfo.GetDesp(), nid)
				return moved, false
			}
			// Move only one node in a loop, we need wait the removing node be removed from raft group.
			break
		}
		if len(namespaceInfo.Removings) > 0 {
			// need wait moved node be removed from raft group
			isAllBalanced = false
			continue
		}
		expectLeader := partitionNodes[namespaceInfo.Partition][0]
		if _, ok := namespaceInfo.Removings[expectLeader]; ok {
			cluster.CoordLog().Infof("namespace %v expected leader: %v is marked as removing", namespaceInfo.GetDesp(),
				expectLeader)
		} else {
			isrList := namespaceInfo.GetISR()
			if len(moveNodes) == 0 && (len(isrList) >= namespaceInfo.Replica) &&
				(namespaceInfo.RaftNodes[0] != expectLeader) {
				for index, nid := range namespaceInfo.RaftNodes {
					if nid == expectLeader {
						cluster.CoordLog().Infof("need move leader for namespace %v since %v not expected leader: %v",
							namespaceInfo.GetDesp(), namespaceInfo.RaftNodes, expectLeader)
						newNodes := make([]string, len(namespaceInfo.RaftNodes))
						copy(newNodes, namespaceInfo.RaftNodes)
						newNodes[0], newNodes[index] = newNodes[index], newNodes[0]
						namespaceInfo.RaftNodes = newNodes
						err := dp.pdCoord.register.UpdateNamespacePartReplicaInfo(namespaceInfo.Name, namespaceInfo.Partition,
							&namespaceInfo.PartitionReplicaInfo, namespaceInfo.PartitionReplicaInfo.Epoch())
						moved = true
						if err != nil {
							cluster.CoordLog().Infof("move leader for namespace %v failed: %v", namespaceInfo.GetDesp(), err)
							return moved, false
						}
					}
				}
				// wait raft leader election
				select {
				case <-monitorChan:
				case <-time.After(time.Second * 5):
				}
			}
		}
		if len(moveNodes) > 0 || moved {
			isAllBalanced = false
		}
	}

	return moved, isAllBalanced
}

type SortableStrings []string

func (s SortableStrings) Less(l, r int) bool {
	return s[l] < s[r]
}
func (s SortableStrings) Len() int {
	return len(s)
}
func (s SortableStrings) Swap(l, r int) {
	s[l], s[r] = s[r], s[l]
}

func getRebalancedNamespacePartitions(ns string,
	partitionNum int, replica int,
	currentNodes map[string]cluster.NodeInfo, balanceVer string) ([][]string, *cluster.CoordErr) {
	if len(currentNodes) < replica {
		return nil, ErrNodeUnavailable
	}
	// for namespace we have much partitions than nodes,
	// so we need make all the nodes have the almost the same leader partitions,
	// and also we need make all the nodes have the almost the same followers.
	// and to avoid the data migration, we should keep the data as much as possible
	// algorithm as below:
	// 1. sort node id ; 2. sort the namespace partitions
	// 3. choose the (leader, follower, follower) for each partition,
	// start from the index of the current node array
	// 4. for next partition, start from the next index of node array.
	//  l -> leader, f-> follower
	// 5. if any node has the replicas more than average, we move the follow from the
	// the most node to the less node(if any conflict choose next less node)
	//         nodeA   nodeB   nodeC   nodeD
	// p1       l       f        f
	// p2               l        f      f
	// p3       f                l      f
	// p4       f       f               l
	// p5       l       f        f
	// p6               l        f      f
	// p7       f                l      f
	// p8       f       f               l
	// p9       l       f        f
	// p10              l        f      f

	// after nodeB is down, the migration as bellow
	//         nodeA   xxxx   nodeC   nodeD
	// p1       l       x        f     x-f
	// p2      x-f      x       f-l     f
	// p3       f                l      f
	// p4       f       x       x-f     l
	// p5       l       x        f     x-f
	// p6      x-f      x       f-l     f
	// p7       f                l      f
	// p8       f       x       x-f     l
	// p9       l       x        f     x-f
	// p10     x-f      x       f-l     f

	// if using the V2 balanced, we choose as below
	// first choose all the replicas one by one (we ignore leader at first step)
	//         nodeA   nodeB   nodeC   nodeD
	// p1       f       f        f
	// p2       f       f               f
	// p3       f                f      f
	// p4               f        f      f
	// p5       f       f        f
	// p6       f       f               f
	// p7       f                f      f
	// p8               f        f      f
	// p9       f       f        f
	// p10      f       f               f
	// then we choose the leader for each partition (first follower which leader is available)
	//         nodeA   nodeB   nodeC   nodeD
	// p1       l       f        f
	// p2       f       l               f
	// p3       f                l      f
	// p4               f        f      l
	// p5       l       f        f
	// p6       f       l               f
	// p7       f                l      f
	// p8               f        f      l
	// p9       l       f        f
	// p10      f       l               f

	// after nodeB is down, the migration as bellow
	//         nodeA   xxxx   nodeC   nodeD
	// p1       l       x        f     x-f
	// p2       f       x       x-l     f
	// p3       f               l-f    f-l
	// p4      x-l      x        f     l-f
	// p5      l-f      x       f-l    x-f
	// p6       f       x       x-f    f-l
	// p7      f-l              l-f     f
	// p8      x-f      x       f-l    l-f
	// p9      l-f      x        f     x-l
	// p10     f-l      x       x-f     f

	// if there are several data centers, we sort them one by one as below
	// nodeA1@dc1 nodeA2@dc2 nodeA3@dc3 nodeB1@dc1 nodeB2@dc2 nodeB3@dc3

	nodeNameList := getNodeNameList(currentNodes)
	return getRebalancedPartitionsFromNameList(ns, partitionNum, replica, nodeNameList, balanceVer)
}

func getRebalancedPartitionsFromNameList(ns string,
	partitionNum int, replica int,
	nodeNameList []SortableStrings, balanceVer string) ([][]string, *cluster.CoordErr) {

	var combined SortableStrings
	sortedNodeNameList := make([]SortableStrings, 0, len(nodeNameList))
	for _, nList := range nodeNameList {
		sortedNodeNameList = append(sortedNodeNameList, nList)
	}
	totalCnt := 0
	for idx, nList := range sortedNodeNameList {
		sort.Sort(nList)
		sortedNodeNameList[idx] = nList
		totalCnt += len(nList)
	}
	if totalCnt < replica {
		return nil, ErrNodeUnavailable
	}
	idx := 0
	for len(combined) < totalCnt {
		nList := sortedNodeNameList[idx%len(sortedNodeNameList)]
		if len(nList) == 0 {
			idx++
			continue
		}
		combined = append(combined, nList[0])
		sortedNodeNameList[idx%len(sortedNodeNameList)] = nList[1:]
		idx++
	}

	if balanceVer == "v2" {
		return fillPartitionMapV2(ns, partitionNum, replica, combined), nil
	}
	if balanceVer == "v3" {
		return fillPartitionMapV3(ns, partitionNum, replica, combined), nil
	}
	return fillPartitionMapV1(ns, partitionNum, replica, combined), nil
}

func fillPartitionMapV1(ns string,
	partitionNum int, replica int,
	sortedNodes SortableStrings) [][]string {

	partitionNodes := make([][]string, partitionNum)
	selectIndex := int(murmur3.Sum32([]byte(ns)))
	for i := 0; i < partitionNum; i++ {
		nlist := make([]string, replica)
		partitionNodes[i] = nlist
		for j := 0; j < replica; j++ {
			nlist[j] = sortedNodes[(selectIndex+j)%len(sortedNodes)]
		}
		selectIndex++
	}
	return partitionNodes
	// maybe check if any node has too much replicas than average
	// to avoid move data, we do not do balance here.(use v2 balance instead if need more balanced)
	//
}

// choose replica by hash, this will cause the leader and replica partly unbalanced, but it will
// reduce the balance data while node failed or added
func fillPartitionMapV2(ns string,
	partitionNum int, replica int,
	sortedNodes SortableStrings) [][]string {

	partitionNodes := make([][]string, partitionNum)
	nameHashMap := make(map[string]uint64, len(sortedNodes))
	for _, nodeName := range sortedNodes {
		h := murmur3.Sum64([]byte(ns + nodeName))
		nameHashMap[nodeName] = h
	}
	//sort hash(hash (node+ns), hash pid)_node
	b := make([]byte, 16)
	for i := 0; i < partitionNum; i++ {
		nlist := make([]string, replica)
		partitionNodes[i] = nlist
		binary.BigEndian.PutUint64(b, uint64(i))
		phash := murmur3.Sum64(b[:8])
		hashSorted := make(SortableStrings, 0, len(sortedNodes))
		for _, nodeName := range sortedNodes {
			nh := nameHashMap[nodeName]
			binary.BigEndian.PutUint64(b, nh)
			binary.BigEndian.PutUint64(b[8:16], phash)
			hashSorted = append(hashSorted, strconv.FormatUint(murmur3.Sum64(b), 10)+"-"+nodeName)
		}
		sort.Sort(hashSorted)
		for j := 0; j < replica; j++ {
			name := hashSorted[j%len(hashSorted)]
			subs := strings.SplitN(name, "-", 2)
			nlist[j] = subs[1]
		}
	}
	return partitionNodes
}

type loadItem struct {
	name        string
	leaderPids  []int
	replicaPids []int
}

func loadItemLeaderCmp(l interface{}, r interface{}) int {
	li := l.(loadItem)
	ri := r.(loadItem)
	if len(li.leaderPids) == len(ri.leaderPids) && len(li.replicaPids) == len(ri.replicaPids) {
		return utils.StringComparator(li.name, ri.name)
	}
	if len(li.leaderPids) == len(ri.leaderPids) {
		return utils.IntComparator(len(li.replicaPids), len(ri.replicaPids))
	}
	return utils.IntComparator(len(li.leaderPids), len(ri.leaderPids))
}

func loadItemReplicaCmp(l interface{}, r interface{}) int {
	li := l.(loadItem)
	ri := r.(loadItem)
	if len(li.replicaPids) == len(ri.replicaPids) {
		return utils.StringComparator(li.name, ri.name)
	}
	return utils.IntComparator(len(li.replicaPids), len(ri.replicaPids))
}

func getMinMaxLoadForLeader(leaders map[string][]int, replicas map[string][]int) (loadItem, loadItem) {
	m := treemap.NewWith(loadItemLeaderCmp)
	for name, lpids := range leaders {
		rpids, _ := replicas[name]
		m.Put(loadItem{
			name:        name,
			leaderPids:  lpids,
			replicaPids: rpids,
		}, name)
	}
	min, _ := m.Min()
	max, _ := m.Max()
	return min.(loadItem), max.(loadItem)
}

func getMinMaxLoadForReplica(replicas map[string][]int) (loadItem, loadItem) {
	m := treemap.NewWith(loadItemReplicaCmp)
	for name, rpids := range replicas {
		m.Put(loadItem{
			name:        name,
			replicaPids: rpids,
		}, name)
	}
	min, _ := m.Min()
	max, _ := m.Max()
	return min.(loadItem), max.(loadItem)
}

func fillPartitionMapV3(ns string,
	partitionNum int, replica int,
	sortedNodes SortableStrings) [][]string {

	oldPartitionNodes := make([][]string, 0)
	newNodesReplicaMap := make(map[string][]int)
	newNodesLeaderMap := make(map[string][]int)
	for _, n := range sortedNodes {
		newNodesReplicaMap[n] = make([]int, 0)
		newNodesLeaderMap[n] = make([]int, 0)
	}
	for pid, olist := range oldPartitionNodes {
		for i, name := range olist {
			if i == 0 {
				pidlist, ok := newNodesLeaderMap[name]
				if ok {
					pidlist = append(pidlist, pid)
					newNodesLeaderMap[name] = pidlist
				}
			}
			pidlist, ok := newNodesReplicaMap[name]
			if ok {
				pidlist = append(pidlist, pid)
				newNodesReplicaMap[name] = pidlist
			}
		}
	}

	partitionNodes := make([][]string, partitionNum)
	// check and reuse the old first and if no old found, we choose the least load on current and update the load
	for pid := 0; pid < partitionNum; pid++ {
		var oldlist []string
		if pid < len(oldPartitionNodes) {
			oldlist = oldPartitionNodes[pid]
		}
		nlist := make([]string, replica)
		partitionNodes[pid] = nlist
		for j := 0; j < replica; j++ {
			var old string
			if len(oldlist) > j {
				old = oldlist[j]
			}
			if j == 0 {
				// check if old leader still alive for leader
				_, ok := newNodesLeaderMap[old]
				if ok {
					nlist[j] = old
					continue
				}
				nleader, _ := getMinMaxLoadForLeader(newNodesLeaderMap, newNodesReplicaMap)
				newNodesLeaderMap[nleader.name] = append(nleader.leaderPids, pid)
				newNodesReplicaMap[nleader.name] = append(nleader.replicaPids, pid)
				nlist[j] = nleader.name
				continue
			}
			_, ok := newNodesReplicaMap[old]
			if ok {
				nlist[j] = old
				continue
			}
			nreplica, _ := getMinMaxLoadForReplica(newNodesReplicaMap)
			newNodesReplicaMap[nreplica.name] = append(nreplica.replicaPids, pid)
			nlist[j] = nreplica.name
		}
	}
	// move if unbalanced
	balanced := false
	maxMoved := replica * partitionNum
	for !balanced {
		partitionNodes, balanced = moveIfUnbalanced(partitionNum, replica, newNodesLeaderMap,
			newNodesReplicaMap, partitionNodes)
		maxMoved--
		if maxMoved < 0 {
			cluster.CoordLog().Warningf("balance moved too much times: %v", partitionNodes)
			break
		}
	}
	return partitionNodes
}

func findPidInList(pid int, l []int) bool {
	for _, p := range l {
		if p == pid {
			return true
		}
	}
	return false
}

// will only move once and return changed replica map and whether we already balanced.
func moveIfUnbalanced(partitionNum int, replica int,
	newNodesLeaderMap map[string][]int, newNodesReplicaMap map[string][]int,
	partitionNodes [][]string) ([][]string, bool) {
	min, max := getMinMaxLoadForLeader(newNodesLeaderMap, newNodesReplicaMap)
	balanced := true
	if len(max.leaderPids)-len(min.leaderPids) <= 1 {
		// leader is balanced
	} else {
		balanced = false
		cluster.CoordLog().Infof("some node has too much leaders: %v", newNodesLeaderMap)
		cluster.CoordLog().Infof("before move, replicas: %v", partitionNodes)
		//
		for _, pid := range max.leaderPids {
			exist := false
			for _, moveTo := range min.leaderPids {
				if pid == moveTo {
					exist = true
					break
				}
			}
			if exist {
				continue
			}
			// pid
			replicas := partitionNodes[pid]
			for i := 0; i < len(replicas); i++ {
				if replicas[i] == max.name {
					// move to min
					replicas[i] = min.name
					break
				}
			}
			min.leaderPids = append(min.leaderPids, pid)
			min.replicaPids = append(min.replicaPids, pid)
			newNodesLeaderMap[min.name] = min.leaderPids
			newNodesReplicaMap[min.name] = min.replicaPids
			npids := make([]int, 0)
			for _, p := range max.leaderPids {
				if p == pid {
					continue
				}
				npids = append(npids, p)
			}
			newNodesLeaderMap[max.name] = npids
			npids = make([]int, 0)
			for _, p := range max.replicaPids {
				if p == pid {
					continue
				}
				npids = append(npids, p)
			}
			newNodesReplicaMap[max.name] = npids
			break
		}
		cluster.CoordLog().Infof("after moved, replicas: %v", partitionNodes)
		return partitionNodes, balanced
	}

	min, max = getMinMaxLoadForReplica(newNodesReplicaMap)
	if len(max.replicaPids)-len(min.replicaPids) <= 1 {
		// replica is balanced
	} else {
		balanced = false
		cluster.CoordLog().Infof("some node has too much replicas: %v", newNodesReplicaMap)
		cluster.CoordLog().Infof("before move, replicas: %v", partitionNodes)
		for _, pid := range max.replicaPids {
			exist := false
			for _, moveTo := range min.replicaPids {
				if pid == moveTo {
					exist = true
					break
				}
			}
			if exist {
				continue
			}
			// pid
			replicas := partitionNodes[pid]
			for i := 0; i < len(replicas); i++ {
				if replicas[i] == max.name {
					// move to min
					replicas[i] = min.name
					break
				}
			}
			min.replicaPids = append(min.replicaPids, pid)
			newNodesReplicaMap[min.name] = min.replicaPids
			npids := make([]int, 0)
			for _, p := range max.replicaPids {
				if p == pid {
					continue
				}
				npids = append(npids, p)
			}
			newNodesReplicaMap[max.name] = npids
			break
		}
		cluster.CoordLog().Infof("after moved, replicas: %v", partitionNodes)
		return partitionNodes, balanced
	}
	return partitionNodes, balanced
}

func (dp *DataPlacement) decideUnwantedRaftNode(namespaceInfo *cluster.PartitionMetaInfo, currentNodes map[string]cluster.NodeInfo) string {
	unwantedNode := ""
	//remove the unwanted node in isr
	partitionNodes, err := getRebalancedNamespacePartitions(
		namespaceInfo.Name,
		namespaceInfo.PartitionNum,
		namespaceInfo.Replica, currentNodes, dp.balanceVer)
	if err != nil {
		return unwantedNode
	}
	for _, nid := range namespaceInfo.GetISR() {
		found := false
		for _, validNode := range partitionNodes[namespaceInfo.Partition] {
			if nid == validNode {
				found = true
				break
			}
		}
		if !found {
			unwantedNode = nid
		}
	}
	return unwantedNode
}
