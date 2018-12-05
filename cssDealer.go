package cssSelector
import (
	"golang.org/x/net/html"
	"strings"
	"regexp"
	"strconv"
	// "fmt"
)

//Given a node and a css query expresion,return a slice of node
func Query(node *html.Node,cssExpression string)(res []*html.Node){
	res=[]*html.Node{}
	//preDeal make it standard Formate
	cssExpression=preDeal(cssExpression)

	//Split into groups and query them
	groups:=getGroups(cssExpression)
	for _,group:=range groups{
		if group!=""{
			res=append(res,QueryFromCssGroup(node,group)...)
		}
	}
	return res
} 
//Given a node and a cssGroup query expression,return a slice of node
func QueryFromCssGroup(node *html.Node,cssGroup string)(res []*html.Node){
	atoms,joiners:=getAtomsAndJoiners(cssGroup)
	rangeSlice:=[]*html.Node{node}
	return queryFromAtomAndJoiners(rangeSlice,atoms,joiners)
}

func queryFromAtomAndJoiners(nodes []*html.Node,atoms,joiners []string)[]*html.Node{
	tmp:=[]*html.Node{}
	if atoms==nil||len(atoms)==0{
		return nodes
	}
	for _,node:=range nodes{
		res:=queryFromCssAtom(node,atoms[0])
		tmp=append(tmp,res...)
	}
	tmp1:=[]*html.Node{}
	if joiners!=nil&&len(joiners)>0{
		for _,n:=range tmp{
			switch(joiners[0]){
				case " ":{
					tmp1=append(tmp1,n)
				}
				case "+":{
					tmp1=append(tmp1,n.NextSibling)
				}
				case ">":{
					tmp1=append(tmp1,ChildElements(n)...)
				}
				case "~":{
					tmp1=append(tmp1,BehindElements(n)...)
				}
			}
		}
		return queryFromAtomAndJoiners(tmp1,atoms[1:],joiners[1:])
	}else{
		return tmp
	}

}

//Deal a css query expression,so that is can be analyzed correctlly
func preDeal(cssExpression string)string{
	fp:=regexp.MustCompile(`\s+`)
	p1:=fp.ReplaceAllLiteralString(cssExpression," ")

	fp1:=regexp.MustCompile(`\s?([+>,~])\s?`)
	p1=fp1.ReplaceAllString(p1,"$1")
	p1=strings.Trim(p1," ")
	return p1
}
//Split a css query expression into several groups
func getGroups(cssStdExpression string)[]string{
	return strings.Split(cssStdExpression,",")
}
//Analyze a css query group，Split it into a slice of css query atom and a slice of joiner
func getAtomsAndJoiners(cssGroupExpression string)(atoms,atomJioners []string){
	// joinMap:=map[string]string{"+":"下一个",">":"子元素"," ":"内的","~","后面的"}
	atomeReg:=regexp.MustCompile(`([+> ~])`)
	atoms=atomeReg.Split(cssGroupExpression,-1)
	atomJioners=atomeReg.FindAllString(cssGroupExpression,-1)
	return 
}
//Given a node and css atom query expression and return a slice of node
func queryFromCssAtom(node *html.Node,cssAtom string)(res []*html.Node){
	res=[]*html.Node{}
	attrReg:=regexp.MustCompile(`^(\w+)?(\[(\w+)([\^\$\*]?=)"([0-9a-zA-Z_-]+)"\])+$`)
	orderReg:=regexp.MustCompile(`^(\w+)?:(((\w+)-)+\w+)\((\d+)\)$`)
	switch{
		case cssAtom=="*":{
			return Traversal(node,AttrFilter{})
		}
		case strings.HasPrefix(cssAtom,"#"):{
			return Traversal(node,AttrFilter{AttrCondis:[]AttrCondtion{AttrCondtion{"id","=",cssAtom[1:]}}})
		}
		case strings.HasPrefix(cssAtom,"."):{
			return Traversal(node,AttrFilter{AttrCondis:[]AttrCondtion{AttrCondtion{"class","*=",cssAtom[1:]}}})
		}
		case strings.HasSuffix(cssAtom,":first-child"):{
			ns:=Traversal(node,OrderFilter{Order:1,PositiveSequense:true})
			for _,n:=range ns{
				if n.DataAtom.String()==strings.TrimSuffix(cssAtom,":first-child"){
					res=append(res,n)
				}
			}
			return res
		}
		case strings.HasSuffix(cssAtom,":last-child"):{
			ns:=Traversal(node,OrderFilter{Order:1})
			for _,n:=range ns{
				if n.DataAtom.String()==strings.TrimSuffix(cssAtom,":last-child"){
					res=append(res,n)
				}
			}
			return res
		}
		case strings.HasSuffix(cssAtom,":first-of-type"):{
			return Traversal(node,OrderFilter{Name: strings.TrimSuffix(cssAtom,":first-of-type"),Order:1,PositiveSequense:true})
		}
		case strings.HasSuffix(cssAtom,":last-of-type"):{
			return Traversal(node,OrderFilter{Name: strings.TrimSuffix(cssAtom,":last-of-type"),Order:1})
		}
		case strings.HasSuffix(cssAtom,":only-child"):{
			if node.NextSibling==nil&&node.NextSibling==nil&&strings.TrimSuffix(cssAtom,":only-child")==node.Data{
				res=append(res,node)
			}
			return res
		}
		case strings.HasSuffix(cssAtom,":empty"):{
			if node.FirstChild==nil&&strings.TrimSuffix(cssAtom,":empty")==node.Data{
				res=append(res,node)
			}
			return res
		}
		case orderReg.MatchString(cssAtom):{
			gs:=orderReg.FindStringSubmatch(cssAtom)
			onum,_:=strconv.Atoi(gs[5])
			onumber:=uint(onum)
			switch gs[2]{
				case "nth-child":{
					ns:=Traversal(node,OrderFilter{Order:onumber,PositiveSequense:true})
					for _,n:=range ns{
						if n.Data==gs[1]{
							res=append(res,n)
						}
					}
					return res
				}
				case "nth-last-child":{
					ns:=Traversal(node,OrderFilter{Order:onumber})
					for _,n:=range ns{
						if n.Data==gs[1]{
							res=append(res,n)
						}
					}
					return res
				}
				case "nth-of-type":{
					return Traversal(node,OrderFilter{Name:gs[1],Order:onumber,PositiveSequense:true})
				}
				case "nth-last-of-type":{
					return Traversal(node,OrderFilter{Name:gs[1],Order:onumber})
				}
			}
			return
		}
		case attrReg.MatchString(cssAtom):{
			tN:=attrReg.FindStringSubmatch(cssAtom)[1]
			cssAtom=cssAtom[len(tN):]
			cssAtom=cssAtom[1:len(cssAtom)-1]
			attrs:=strings.Split(cssAtom,"][")
			attrFilter:=AttrFilter{Name:tN,AttrCondis:[]AttrCondtion{}}
			attrExpReg:=regexp.MustCompile(`^(\w+)(([\^\$\*]?=)"([0-9a-zA-Z_-]+)")?$`)
			for _,attr:=range attrs{
				if attrExpReg.MatchString(attr){
					gs:=attrExpReg.FindStringSubmatch(attr)
					attrFilter.AttrCondis=append(attrFilter.AttrCondis,AttrCondtion{Name:gs[1],Operator:gs[3],Value:gs[4]})
				}
			}
			return Traversal(node,attrFilter)
		}
		default:{
			return Traversal(node,AttrFilter{Name:cssAtom})
		}
	}
}

//Given a node and Filter implementation and return a slice of node.
func Traversal(node *html.Node,filter Filter)(res []*html.Node){
	res=[]*html.Node{}
	if node==nil{
		return 
	}
	if filter.Accept(node){
		res=append(res,node)
	}
	firstChild:=node.FirstChild
	if firstChild!=nil{
		res=append(res,Traversal(firstChild,filter)...)
		for next:=firstChild.NextSibling;next!=nil;next=next.NextSibling{
			res=append(res,Traversal(next,filter)...)
		}
	}
	return res
}

type AttrFilter struct{
	Name string
	AttrCondis []AttrCondtion
}
func (af AttrFilter)Accept(node *html.Node)bool{
	if node==nil{
		return false
	}
	if af.Name!=""&&node.DataAtom.String()!=af.Name{
		return false
	}
	attrMap:=Attribute2Map(node.Attr)
	if af.AttrCondis!=nil{
		for _,attrCondi:=range af.AttrCondis{
			if val,ok:=attrMap[attrCondi.Name];!ok{
				return false
			}else{
				switch(attrCondi.Operator){
					case "^=":{
						if !strings.HasPrefix(val,attrCondi.Value){
							return false
						}
					}
					case "$=":{
						if !strings.HasSuffix(val,attrCondi.Value){
							return false
						}
					}
					case "*=":{
						if !strings.Contains(val,attrCondi.Value){
							return false
						}
					}
					case "=":{
						if val!=attrCondi.Value{
							return false
						}
					}
				}
			}
		}
	}
	return true
}
type OrderFilter struct{
	Name string
	Order uint
	PositiveSequense bool
}
func (of OrderFilter)Accept(node *html.Node)bool{
	if node.Type!=html.ElementNode{
		return false
	}
	var o uint=1
	if of.PositiveSequense{
		if of.Name!=""{
			for prev:=node.PrevSibling;prev!=nil;prev=prev.PrevSibling{
				if prev.Type==html.ElementNode&&node.DataAtom.String()==of.Name{
					o++
				}
			}
		}else{
			for prev:=node.PrevSibling;prev!=nil;prev=prev.PrevSibling{
				if prev.Type==html.ElementNode{
					o++
				}
			}
		}
	}else{
		if of.Name!=""{
			for next:=node.NextSibling;next!=nil;next=next.NextSibling{
				if next.Type==html.ElementNode&&node.DataAtom.String()==of.Name{
					o++
				}
			}
		}else{
			for next:=node.NextSibling;next!=nil;next=next.NextSibling{
				if next.Type==html.ElementNode{
					o++
				}
			}
		}
	}
	return o==of.Order
}
type AttrCondtion struct{
	Name ,Operator ,Value string
}

type Filter interface{
	Accept(node *html.Node)bool
}

func Attribute2Map(attrs []html.Attribute)map[string]string{
	attrMap:=map[string]string{}
	if attrs!=nil{
		for _,attr:=range attrs{
			attrMap[attr.Key]=attr.Val
		}
	}
	return attrMap
}


//Get all child node of the param node
func ChildElements(node *html.Node)[]*html.Node{
	res:=[]*html.Node{}
	if node.FirstChild!=nil{
		for next:=node.FirstChild;next!=nil;next=next.NextSibling{
			if next.Type==html.ElementNode{
				res=append(res,next)
			}
		}
	}
	return res
}
//Get all elements node behind the param node
func BehindElements(node *html.Node)[]*html.Node{
	res:=[]*html.Node{}
	for next:=node.NextSibling;next!=nil;next=next.NextSibling{
		if next.Type==html.ElementNode{
			res=append(res,next)
		}
	}
	return res
}

func InnerText(node *html.Node)string{
	text:=""
	for n:=node.FirstChild;n!=nil;n=n.NextSibling{
		if n.Type==html.TextNode{
			text+=" "+n.Data
		}else{
			text+=" "+InnerText(n)
		}
	}
	return text
}