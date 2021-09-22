package main

import (
	"bufio"
	"fmt"
	bitset "github.com/bits-and-blooms/bitset"
	"math"
	"os"
	"strconv"
	"time"
)

type  Node struct{
	m_low int
	m_hight int
	m_content string
	m_type int
}

func newNode(m_low,m_hight int,m_content string,m_type int) Node {
	node := Node{
		m_low: m_low,
		m_hight: m_hight,
		m_content: m_content,
		m_type: m_type,
	}
	return node
}

func bitCompute (bstr string,size int) bitset.BitSet{
	bint, _ :=strconv.Atoi(bstr)
	bturn := bint << size
	buf:=[]uint64{uint64(bturn)}
	newb := bitset.From(buf)
	return *newb
}

type delta []Node

func backtrackingPath(basic string,s string) delta {
	m := len(basic)
	n := len(s)
	dp := make([][]int,m+1)
	for i:=0;i<m+1;i++ {
		dp[i] = make([]int,n+1)
	}
	for i := 0;i<m+1;i++ {
		dp[i][0]=i
	}
	for i := 0;i<n+1;i++ {
		dp[0][i]=i
	}
	for i := 1;i < m+1;i++ {
		for j := 1; j < n+1; j++ {
			if basic[i-1] == s[j-1] {
				dp[i][j] = dp[i-1][j-1]
			}else{
				dp[i][j] = int(math.Min(float64(dp[i-1][j-1]+1),math.Min(float64(dp[i-1][j]+1),float64(dp[i][j-1]+1))))
			}
		}
	}
	fmt.Printf("distance %n",dp[m][n])
	stack_delta := make([]Node,0)
	for n>=0||m>=0 {
		if n!=0&&dp[m][n-1]+1 ==dp[m][n] {
			node := newNode(m-1,m, string(s[n-1]),0)
			if len(stack_delta)==0 {
				stack_delta = append(stack_delta, node)
			}else{
				//tempNode := stack_delta[len(stack_delta)-1:len(stack_delta)]
				//这样写返回的是一个切片数组,tempNode[0]是对象
				tempNode := stack_delta[len(stack_delta)-1]
				if tempNode.m_low==node.m_low&&tempNode.m_type==node.m_type{
					node.m_content=node.m_content+tempNode.m_content
					stack_delta = stack_delta[:len(stack_delta)-1]
				}
				stack_delta = append(stack_delta, node);
			}
			fmt.Printf("insert %s  at %d \n",s[n-1],m-1)
			n = n-1
			continue
		} else if m!=0 &&dp[m-1][n]+1 ==dp[m][n]{
			fmt.Printf("delete %s  at  %d \n",basic[m-1],m-1)
			node := newNode(m-1,m-1,"-",1)
			if len(stack_delta)==0 {
				stack_delta = append(stack_delta, node)
			}else{
				//tempNode := stack_delta[len(stack_delta)-1:len(stack_delta)]
				//这样写返回的是一个切片数组,tempNode[0]是对象
				tempNode := stack_delta[len(stack_delta)-1]
				if tempNode.m_low==node.m_low+1&&tempNode.m_type==node.m_type{
					node.m_hight=tempNode.m_hight
					stack_delta = stack_delta[:len(stack_delta)-1]
				}
				stack_delta = append(stack_delta, node);
			}
			m = m-1
			continue
		} else if (m>=1&&n>=1&&dp[m-1][n-1]+1 == dp[m][n]){
			fmt.Printf("replace %s  to  %s  at  %d \n",basic[m-1],s[n-1],m-1)
			node := newNode(m-1,m-1,string(s[n-1]),2)
			if len(stack_delta)==0 {
				stack_delta = append(stack_delta, node)
			}else{
				//tempNode := stack_delta[len(stack_delta)-1:len(stack_delta)]
				//这样写返回的是一个切片数组,tempNode[0]是对象
				tempNode := stack_delta[len(stack_delta)-1]
				if tempNode.m_low==node.m_low+1 && tempNode.m_type==node.m_type{
					node.m_hight=tempNode.m_hight
					node.m_content = node.m_content + tempNode.m_content
					stack_delta = stack_delta[:len(stack_delta)-1]
				}
				stack_delta = append(stack_delta, node);
			}
			n = n-1
			m = m-1
			continue
		}
		n = n - 1
		m = m - 1
	}
	return stack_delta
}

func main()  {
	startTime := time.Now()
	fmt.Println("deltaProject is run.")

	//输出字符串
	foutpath := "./resource/outfile.txt"
	outfile,outErr := os.Create(foutpath)
	if outErr !=nil{
		fmt.Println(outErr.Error())
	}
	defer outfile.Close()

	//基础字符串
	basicpath :="./resource/chr0716.txt"
	basic_file,basic_err := os.Open(basicpath)
	if basic_err !=nil{
		panic(basic_err)//当系统发现无法继续运行下去的故障时调用，导致程序终止，然后由系统显示错误号
	}
	defer basic_file.Close()
	basic_reader := bufio.NewReader(basic_file)
	basicStr,basic_err := basic_reader.ReadString('\n')

	//原始数据
	filename := "./resource/query.txt"
	file,err := os.Open(filename)
	if err !=nil{
		panic(err)//当系统发现无法继续运行下去的故障时调用，导致程序终止，然后由系统显示错误号
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		stackdelta := backtrackingPath(basicStr,scanner.Text())
		for len(stackdelta)!=0 {
			stackNode := stackdelta[len(stackdelta)-1]
			//bitset操作
			bit_ := bitset.New(24)
			bit_.Set(uint(stackNode.m_low))
			bstr := bit_.String()
			basic_b := bitCompute(bstr,11)
			hb :=bitset.New(24)
			hb.Set(uint(stackNode.m_hight))
			basic_bb := basic_b.SymmetricDifference(hb)
			basic_bbb := bitCompute(basic_bb.String(),2)
			tb :=bitset.New(24)
			tb.Set(uint(stackNode.m_type))
			final_b := basic_bbb.SymmetricDifference(tb)
			outfile.WriteString(final_b.String())
			content_size := len(stackNode.m_content)
			bit_content := bitset.New(16)
			bit_content.Set(uint(content_size))
			outfile.WriteString(bit_content.String())
			outfile.WriteString(stackNode.m_content+"\n")
			stackdelta = stackdelta[:len(stackdelta)-1]
		}
	}
	if err := scanner.Err();err!=nil {
		panic(err)
	}
	endTime := time.Since(startTime)//程序所用时间
	fmt.Println(endTime)
}
