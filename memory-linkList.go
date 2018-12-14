package mem
//
//import (
//	"Platon-go/log"
//	"fmt"
//)
//
//var (
//	chunkBaseSize = 1
//)
//
//const (
//	Malloc    = iota // value --> 0
//	MSuitable        // value --> 1
//	MGetSpilt
//	MSpilt  // value --> 2
//	MMerge
//
//	Free
//	FSuitable
//	FSpilt
//	FSpiltRecover
//	FMerge
//	Rest
//)
//
//type Statistic struct {
//	Malloc        int
//	MSuitable     float32
//	MGetSpilt     float32
//	MSpilt        float32
//	MMerge        float32
//	Free          int
//	FSuitable     float32
//	FSpilt        float32
//	FSpiltRecover float32
//	FMerge        float32
//	Rest          float32
//}
//
//type Memory struct {
//	Memory       []byte
//	Start        int
//	Size         int
//	MaxSize      int
//	Deep         int //最大深度
//	MemTables    []*TableHead
//	Inited       bool
//	BigMemLL     *expendHead
//	BigMemFreeLL *expendHead
//	extend       map[int]int //扩展内存的起止地址
//}
//
//type BigMem struct {
//	start, end int
//	dirty      int
//	head       *ChunkHead
//}
//
//type TableHead struct {
//	FreeNum int //内存块数目
//	Size    int //所属于的memTable
//	Head    *ChunkHead
//	Used    map[int]*Chunk
//	Spilt   map[*Chunk]*Chunk
//	Statis  *Statistic
//}
//
//type Chunk struct {
//	//table int //所属的memTable
//	addr  int //起始位置
//	Size  int //数据占用的大小 :size为0表示未有使用，不为0表示正在使用
//	dirty int //是否被占用，0 没有，1 被占用
//	head  *ChunkHead
//}
//
//type ChunkHead struct {
//	owner *Chunk
//	next  *ChunkHead
//	prev  *ChunkHead
//}
//
//type expendHead struct {
//	start int
//	Size  int
//	page  int
//	head  *ChunkHead
//}
//
////测试所有free之后链表得连续性和数量是否正确
//func (table *TableHead) CheckContinuity() (bool, int) {
//	i := 1
//	for pos := table.Head.next; pos != table.Head.prev; pos = pos.next {
//		if pos.next.owner.addr-pos.owner.addr != table.Size {
//			return false, i
//		}
//		i++
//	}
//	return true, i
//}
//
////链表初始化
//func (table *TableHead) initTableHead() {
//	table.Head.next = table.Head
//	table.Head.prev = table.Head
//}
//
////插入结点
//func addChunk(new *ChunkHead, prev *ChunkHead, next *ChunkHead) {
//	next.prev = new
//	new.next = next
//	new.prev = prev
//	prev.next = new
//}
//
////在链表头部插入
//func addChunkFromHead(new *ChunkHead, head *ChunkHead) {
//	if head == nil || (head.prev == nil && head.next == nil) {
//		new.prev = new
//		new.next = new
//		return
//	}
//	addChunk(new, head, head.next)
//}
//
//func addChunkBehind(new *ChunkHead, head *ChunkHead) {
//	if head == nil || (head.prev == nil && head.next == nil) {
//		new.prev = new
//		new.next = new
//		return
//	}
//	addChunk(new, head.prev, head)
//}
//
//func deleteChunk(entry *ChunkHead) {
//	//entry.next.prev = entry.prev
//	//entry.prev.next = entry.next
//	deleteChunks(entry, entry)
//}
//
//func deleteChunks(start *ChunkHead, end *ChunkHead) {
//	end.next.prev = start.prev
//	start.prev.next = end.next
//	start.prev = nil
//	end.next = nil
//}
//
//func recoverChunk(entry, head *ChunkHead) {
//	recoverChunks(entry, entry, head)
//}
//
//func recoverChunks(start, end, head *ChunkHead) {
//	prev, next := getPrevAndNext(start, head)
//	start.prev = prev
//	end.next = next
//	next.prev = end
//	prev.next = start
//}
//
///**
//查找start所在的前后chunk
//*/
//func getPrevAndNext(start, head *ChunkHead) (*ChunkHead, *ChunkHead) {
//	if head.next == head {
//		//所有chunk已用完
//		return head, head
//	}
//	if head.next.owner.addr > start.owner.addr {
//		//在头部插入
//		return head, head.next
//	}
//
//	if head.prev.owner.addr < start.owner.addr {
//		//在尾部插入
//		return head.prev, head
//	}
//
//	for pos := head.next; pos != head; pos = pos.next {
//		//非最后一个块
//		if pos.owner.addr < start.owner.addr && pos.next.owner.addr > start.owner.addr {
//			return pos, pos.next
//		}
//	}
//	return nil, nil
//}
//
////将内存进行切分成不同的分配阶table
//func (m *Memory) Init() int {
//
//	// 以1byte进行分配
//
//	for i := 0; i < m.Deep; i++ {
//		//初始化memTables
//		statis := Statistic{
//			Malloc:        0,
//			MSuitable:     0,
//			MGetSpilt:     0,
//			MSpilt:        0,
//			MMerge:        0,
//			Free:          0,
//			FSuitable:     0,
//			FSpilt:        0,
//			FSpiltRecover: 0,
//			FMerge:        0,
//			Rest:          0,
//		}
//		m.MemTables[i] = &TableHead{
//			FreeNum: 0,
//			Size:    chunkBaseSize * 1 << uint(i),
//			Head:    &ChunkHead{},
//			Used:    make(map[int]*Chunk),
//			Spilt:   make(map[*Chunk]*Chunk),
//			Statis:  &statis,
//		}
//		m.MemTables[i].initTableHead()
//
//	}
//
//	pmem := m.Start //总内存的首地址
//	m.init(chunkBaseSize, pmem)
//
//	m.Inited = true
//
//	return 0
//}
//
//func (m *Memory) init(chunkBaseSize, pmem int) {
//	tableSize := m.Size / m.Deep
//
//	for i := 0; i < m.Deep; i++ {
//		baseSize := chunkBaseSize << uint(i)
//		num := tableSize / baseSize
//
//		pmem = buildChain(num, pmem, baseSize, tableSize, m.MemTables[i].Head)
//
//		m.MemTables[i].FreeNum = num
//	}
//}
//
//func buildChain(num, pmem, baseSize, tableSize int, head *ChunkHead) int {
//
//	for j := 0; j < num; j++ {
//		tmpChunk := Chunk{
//			addr:  pmem,      //起始位置
//			Size:  baseSize, //数据占用的大小
//			dirty: 0,         //初始化内存卡未被占用
//			head:  &ChunkHead{},
//		}
//		tmpChunk.head.owner = &tmpChunk
//
//		addChunkBehind(tmpChunk.head, head) //插入链表
//		pmem += baseSize                   //指针相应往后移动chunkSize
//	}
//	return pmem
//
//}
//
//func (table *TableHead) getAvailableChunk() *Chunk {
//	if table.FreeNum == 0 {
//		return nil
//	}
//	for pos := table.Head.next; pos != table.Head; pos = pos.next {
//		if pos.owner.dirty == 0 {
//			pos.owner.dirty = 1
//			return pos.owner
//		}
//	}
//
//	return nil
//}
//
////根据size获取相应分配阶在memTables中的index
//func (m Memory) getTableIndex(Size int) int {
//	var tmp int
//	Size /= chunkBaseSize
//
//	for i := 0; tmp < Size; i++ {
//		tmp = 1 << uint(i)
//		if tmp == Size {
//			return i
//		}
//	}
//	return -1 //实际上经过前面判断不可能执行到这
//}
//
//func (m Memory) getTableIndexFromOffset(offset int) int {
//
//	index := (offset - m.Start) / (m.Size / m.Deep)
//	if index > 7 {
//		panic(fmt.Sprintf("invalid offset:%d :index:%d", offset, index))
//	}
//	return index //实际上经过前面判断不可能执行到这
//}
//
//func (m *Memory) fixSize(Size int) int {
//
//	result := 1
//	for result < Size {
//		result = result << 1
//	}
//	return result
//}
//
////初始化扩展链表
//func initExpendLL(m *Memory) {
//	if m.BigMemLL == nil && m.BigMemFreeLL == nil {
//		m.BigMemLL = &expendHead{
//			start: len(m.Memory),
//			head: &ChunkHead{
//			},
//		}
//		m.BigMemLL.head.next = m.BigMemLL.head
//		m.BigMemLL.head.prev = m.BigMemLL.head
//
//		m.BigMemFreeLL = &expendHead{
//			start: len(m.Memory),
//			head: &ChunkHead{
//			},
//		}
//		m.BigMemFreeLL.head.next = m.BigMemFreeLL.head
//		m.BigMemFreeLL.head.prev = m.BigMemFreeLL.head
//	}
//
//}
//
////expend 02 ,add expendHead
//func expend(Size int, m *Memory) int {
//	pos := expendAvailable(m, Size)
//	if pos != -1 {
//		return pos
//	}
//
//	initExpendLL(m)
//
//	pos = len(m.Memory)
//	//扩展内存，1页
//	m.extend[pos] = pos + Size //可以去掉
//
//	m.Memory = append(m.Memory, make([]byte, Size)...)
//
//	chunk := Chunk{
//		addr:  pos,
//		Size:  Size,
//		dirty: 1,
//	}
//	chunk.head = &ChunkHead{
//		owner: &chunk,
//	}
//	addChunkFromHead(chunk.head, nil)
//	addChunkBehind(chunk.head, m.BigMemLL.head)
//
//	return pos
//}
//
//func expendFree(offset int, m *Memory) {
//	for posHead := m.BigMemLL.head.next; posHead != m.BigMemLL.head; posHead = posHead.next {
//		if posHead.owner.addr == offset {
//			deleteChunk(posHead)
//			posHead.owner.dirty = 0
//			//放入freeLL
//			addChunkToFreeLLAndSort(posHead, m.BigMemFreeLL.head)
//
//		}
//	}
//
//}
//
//func expendAvailable(m *Memory, Size int) int {
//
//	return -1
//}
//
////expend 01
////func (m *Memory) expend(Size int) int{
////	pos := expendAvailable(m, Size)
////	if pos != -1 {
////		return pos
////	}
////
////	defaultPageSize :=1024
////	pos = len(m.Memory)
////	//扩展内存，1页
////	m.extend[pos] = pos+ defaultPageSize
////	m.Memory = append(m.Memory, make([]byte, defaultPageSize)...)
////	chunk := Chunk{
////		addr:pos ,
////		Size:Size,
////		dirty:1,
////	}
////	chunk.head = &ChunkHead{
////		owner:&chunk,
////	}
////	addChunkFromHead(chunk.head,nil)
////	m.BigMemLL= append(m.BigMemLL,chunk.head)
////
////	chunkFree:= Chunk{
////		addr:  pos + Size,
////		Size:  defaultPageSize-Size,
////		dirty: 0,
////	}
////	chunkFree.head = &ChunkHead{
////		owner:&chunkFree,
////	}
////	addChunkFromHead(chunkFree.head,nil)
////	m.BigMemFreeLL = append(m.BigMemFreeLL,chunkFree.head)
////	return pos
////}
//
//func addChunkToFreeLLAndSort(new *ChunkHead, ll *ChunkHead) {
//	pos := ll
//	for {
//		if pos.next.owner.addr > new.owner.addr && pos.owner.addr < new.owner.addr {
//			//todo 检查是否存在可以merge的内存块
//			checkMerge(pos, pos.next, pos.prev)
//
//			addChunk(new, pos, pos.next)
//			break
//		}
//		if pos.next == ll {
//			break
//		}
//		pos = pos.next
//	}
//
//}
//
//func checkMerge(pos *ChunkHead, next *ChunkHead, prev *ChunkHead) *ChunkHead {
//	if next != nil && next.owner != nil && pos.owner.addr+pos.owner.Size == next.owner.addr {
//		//合并下一个块
//		chunkMerge(pos, next)
//	}
//	if prev != nil && prev.owner != nil && prev.owner.addr+prev.owner.Size == pos.owner.addr {
//		//合并上一个块
//		chunkMerge(prev, pos)
//	}
//	return prev
//}
//
//func chunkMerge(prev *ChunkHead, next *ChunkHead) {
//	prev.next = next.next
//	next.prev = prev
//
//	prev.owner.Size += next.owner.Size
//}
//
//func (m *Memory) Malloc(Size int) (int) {
//
//	fixSize := m.fixSize(Size)
//
//	if fixSize > 1<<uint(m.Deep-1) {
//		//超过最大分配内存128,从大内存开始往下合并
//		pos := getFromMerge(m.Deep-1, m, Size)
//		if pos == -1 {
//			//pos = m.expend(Size)
//		}
//		return pos
//	}
//
//	index := m.getTableIndex(fixSize)
//	table := m.MemTables[index]
//
//	//从最合适内存获取
//	ret := getFromSuitable(table)
//	if ret != -1 {
//		memStatis(table, Malloc)
//		memStatis(table, MSuitable)
//		memStatis(table, Rest)
//		return ret
//	}
//
//	//从大内存拆分获取
//	ret = getFromUpper(index, m)
//	if ret != -1 {
//		return ret
//	}
//
//	// 从小内存合并
//	ret = getFromMerge(index-1, m, Size)
//	if ret != -1 {
//		return ret
//	}
//	panic("malloc error...")
//	return -1
//}
//
///**
//  小内存合并为大内存,size为所需要的大小
//*/
//func getFromMerge(index int, m *Memory, Size int) int {
//
//	num := 0         //所需要的块数
//	var start *Chunk //起始块
//	var end *Chunk   //末尾块
//	for i := index; i > -1; i-- {
//		num = Size / m.MemTables[i].Size
//		if Size%m.MemTables[i].Size > 0 {
//			num += 1
//		}
//
//		start, end = merge(m.MemTables[i], num)
//		if start != nil && end != nil {
//			//将找到的start至end集体断链
//			deleteChunks(start.head, end.head)
//
//			//将start放入used
//			m.MemTables[i].Used[start.addr] = start
//
//			m.MemTables[i].FreeNum -= num
//
//			memStatis(m.MemTables[i], Malloc)
//			memStatis(m.MemTables[i], MMerge)
//			memStatis(m.MemTables[i], Rest)
//
//			//fmt.Printf("get from merge start:%d,end:%d,num:%d \n", start.addr, end.addr, num)
//			return start.addr
//		}
//
//	}
//
//	return -1
//
//}
//
//func merge(table *TableHead, num int) (*Chunk, *Chunk) {
//	if table == nil || num == 0 || num > table.FreeNum {
//		return nil, nil
//	}
//	//寻找到符合的地址连续的chunk
//	i := 0
//	var start *Chunk
//	var end *Chunk
//	for pos := table.Head.next; pos != table.Head; pos = pos.next {
//		if pos.next == table.Head || pos.next.owner.addr-pos.owner.addr == table.Size {
//			i++
//			if i == 1 {
//				start = pos.owner
//			}
//			if i == num {
//				end = pos.owner
//				break
//			}
//		} else {
//			i = 0
//			start = nil
//		}
//	}
//
//	if start == nil || end == nil {
//		return nil, nil
//	}
//
//	//判断找到的连续的内存块是否正确
//	if end.addr-start.addr != table.Size*(num-1) {
//		fmt.Printf("error get start:%d ,end:%d,needSize:%d \n", start.addr, end.addr, table.Size*num)
//		return nil, nil
//	}
//
//	return start, end
//}
//
////从最合适的table中获取内存
//func getFromSuitable(table *TableHead) int {
//	temp := table.getAvailableChunk()
//	if temp != nil {
//		//fmt.Printf("get from suitable,addr = %d\n", temp.addr)
//
//		//该chunk断链
//		deleteChunk(temp.head)
//		//used里面增加一条数据
//		table.Used[temp.addr] = temp
//		table.FreeNum--
//
//		return temp.addr
//	} else {
//		return -1
//	}
//}
//
////从大内存中获取拆分之后的合适内存
//func getFromUpper(index int, m *Memory) int {
//	//我们从比它大的内存块中再去查找
//
//	for i := index + 1; i < m.Deep; i++ {
//
//		//查找已经做过拆分的块中是否有数据
//		temp := getFromSpilt(m.MemTables[i], m.MemTables[index].Size)
//		if temp != nil {
//			temp.dirty = 1
//			m.MemTables[i].Used[temp.addr] = temp
//
//			memStatis(m.MemTables[i], Malloc)
//			memStatis(m.MemTables[i], MGetSpilt)
//			memStatis(m.MemTables[i], Rest)
//
//			fmt.Printf("get from Spilt index=%d ,addr=%d\n", i, temp.addr)
//			return temp.addr
//		}
//
//		//获取新的块做拆分
//		spiltChunk := m.MemTables[i].getAvailableChunk()
//		if spiltChunk != nil {
//
//			temp = splitChunk(index, i, spiltChunk, m)
//			if temp != nil {
//
//				temp.dirty = 1
//				//放入used
//				m.MemTables[i].Used[temp.addr] = temp
//				m.MemTables[i].Spilt[spiltChunk] = temp
//
//				memStatis(m.MemTables[i], Malloc)
//				memStatis(m.MemTables[i], MSpilt)
//				memStatis(m.MemTables[i], Rest)
//
//				//fmt.Printf("Spilt chunk index=%d,addr:%d \n", i, temp.addr)
//				return temp.addr
//			}
//		}
//
//	}
//	return -1
//}
//
///**
//  从已经拆分的chunk中查找符合size大小的内存
//*/
//func getFromSpilt(table *TableHead, Size int) *Chunk {
//	if len(table.Spilt) == 0 {
//		return nil
//	}
//	for k, v := range table.Spilt {
//		if v.Size != Size {
//			//判断已经拆分的chunk,size是否匹配
//			continue
//		}
//		pos := v.head
//		for i := 0; i < k.Size/v.Size; i++ {
//			if pos.owner.dirty == 0 {
//				return pos.owner
//			}
//			pos = pos.next
//		}
//	}
//	return nil
//}
//
///**
//  拆分chunk,返回拆分之后新的双向循环链表，
//*/
//func splitChunk(dstIndex int, srcIndex int, chunk *Chunk, m *Memory) *Chunk {
//	chunkNum := 1 << (uint(srcIndex - dstIndex)) //2^差值 倍
//
//	srcTable := m.MemTables[srcIndex]
//	deleteChunk(chunk.head) //把被拆分的大内存块从相应链表执行上断链
//
//	spilt := spilt(chunk, chunkNum) //返回拆分后的chunk链表
//	srcTable.FreeNum--              //内存块数目-1
//
//	return spilt
//}
//
///**
//将某一个chunk 等分times次，然后返回chunk的链表
//*/
//func spilt(chunk *Chunk, times int) *Chunk {
//	pmem := chunk.addr
//	dstSize := chunk.Size / times
//	var tempChunk = &Chunk{}
//
//	pre := &ChunkHead{}
//
//	for i := 1; i <= times; i++ {
//		tempChunk = &Chunk{
//			addr:  pmem,
//			Size:  dstSize,
//			dirty: 0,
//			head:  &ChunkHead{},
//		}
//		tempChunk.head.owner = tempChunk
//		addChunkBehind(tempChunk.head, pre)
//		pre = tempChunk.head
//		pmem += dstSize
//	}
//
//	return pre.owner
//}
//
///**
//
// */
//func (m *Memory) Free(offset int) int {
//
//	index := m.getTableIndexFromOffset(offset)
//	//fmt.Printf("Free %d from index=%d \n", offset, index)
//	table := m.MemTables[index]
//
//	defer memStatis(table, Free)
//	defer memStatis(table, Rest)
//
//	//去当前table的used中查找是否有占用
//	if len(table.Used) == 0 {
//		log.Warn(fmt.Sprintf("Free failed , Used.len = 0 \n"))
//		return 0
//	}
//	freeChunk := table.Used[offset]
//	if freeChunk == nil {
//		log.Warn(fmt.Sprintf("can not found Free offset in Used, offset = %d \n", offset))
//		return 0
//	}
//	if freeChunk.addr != offset {
//		log.Error(fmt.Sprintf("freeChunk.addr != offset, offset = %d,addr=%d \n", offset, freeChunk.addr))
//		return -1
//	}
//
//	//清空数据
//
//	//copy(m.Memory[offset:offset+freeChunk.Size], make([]byte, freeChunk.Size))
//	clear(offset, offset+freeChunk.Size, m.Memory)
//
//	//删除对应used记录
//	delete(table.Used, offset)
//
//	if freeChunk.Size == table.Size {
//		//没有做过拆分
//		var end *Chunk
//		i := 0
//
//		pos := freeChunk.head
//		for {
//			i++
//			end = pos.owner
//			pos.owner.dirty = 0
//			pos = pos.next
//			if pos == nil {
//				break
//			}
//		}
//
//		//将对应的chunk恢复至原来链中
//		recoverChunks(freeChunk.head, end.head, table.Head)
//		table.FreeNum += i
//		//fmt.Printf("Free from Used,start:%d,end:%d,num:%d \n",freeChunk.addr,end.addr,i)
//
//		if i == 1 {
//			memStatis(table, FSuitable)
//		} else {
//			memStatis(table, FMerge)
//		}
//
//	} else {
//		//做过拆分，将chunk对应spilt中的数据恢复
//		freeChunk.dirty = 0
//		//判断其他拆分的出来的chunk时候有被占用，都没有占用就将被拆分的chunk恢复至原来链中，有一个占用都不恢复
//		num := table.Size / freeChunk.Size
//		pos := freeChunk.head
//		recover := true
//		for i := 0; i < num; i++ {
//			if pos.owner.dirty == 1 {
//				recover = false
//				break
//			}
//			pos = pos.next
//		}
//
//		//fmt.Printf("Free from Spilt,addr:%d \n",offset)
//
//		memStatis(table, FSpilt)
//
//		if recover {
//			for k := range table.Spilt {
//				if k.addr <= offset && k.addr+k.Size > offset {
//					recoverChunk(k.head, table.Head) //恢复至原来链中
//					delete(table.Spilt, k)           //清空spilt数据
//					k.dirty = 1
//					table.FreeNum++
//					memStatis(table, FSpiltRecover)
//					//fmt.Printf("recover Spilt to table ,addr:%d,,, \n",offset)
//					break
//				}
//			}
//		}
//	}
//
//	return 0
//
//}
//
//func clear(start, end int, mem []byte) {
//	for i := start; i < end; i++ {
//		mem[i] = 0
//	}
//}
//
//func memStatis(head *TableHead, t int) {
//	switch t {
//	case Malloc:
//		head.Statis.Malloc++
//	case MSuitable:
//		head.Statis.MSuitable++
//	case MGetSpilt:
//		head.Statis.MGetSpilt++
//	case MSpilt:
//		head.Statis.MSpilt++
//	case MMerge:
//		head.Statis.MMerge++
//	case Free:
//		head.Statis.Free++
//	case FSuitable:
//		head.Statis.FSuitable++
//	case FSpiltRecover:
//		head.Statis.FSpiltRecover++
//	case FSpilt:
//		head.Statis.FSpilt++
//	case FMerge:
//		head.Statis.FMerge++
//	case Rest:
//		head.Statis.Rest = float32(head.FreeNum)
//
//	}
//
//}
