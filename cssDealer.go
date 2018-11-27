package cssSelector
import (
	"golang.org/x/net/http"
	"strings"
	"regexp"
)

//Given a node and a css query expresion,return a slice of node
func QueryFromCssExpression(node *html.Node,cssExpression string)(res []*html.Node){
	//preDeal make it standard Formate
	cssExpression=preDeal(cssExpression)

	//Split into groups and query them
	groups:=getGroups(cssExpression)
	for _,group:=range groups{
		res=append(res,QueryFromCssGroup(node,group))
	}
	return res
} 
//Given a node and a cssGroup query expression,return a slice of node
func QueryFromCssGroup(node *html.Node,cssGroup string)(res []*html.Node){
	atoms,joiners:=getAtomsAndJoiners(group)
	tmp:=[]*html.Node
	for i:=0;i<len(atoms);i++{
		atomReses:=queryFromCssAtom(atoms[i])
		for _,atomRes:=range atomReses{
			switch(joiners[i]){
				case " ":{
					tmp=append(tmp,Traversal(atomRes,AttrFilter{})...)
				}
				case "+":{
					tmp=append(tmp,atomRes.NextSibling)
				}
				case ">":{
					tmp=append(tmp,ChildElements(atomRes)...)
				}
				case "~":{
					tmp=append(tmp,BehindElements(atomRes)...)
				}
			}
		}
	}
}
//Deal a css query expression,so that is can be analyzed correctlly
func preDeal(cssExpression string)string{
	fp:=regexp.MustCompile(`\s+`)
	p1:=fp.ReplaceAllLiteralString(pattern," ")

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
	atoms:=atomeReg.Split(itemSelector,-1)
	atomJioners:=atomeReg.FindAllString(itemSelector,-1)
	return 
}
//Given a node and css atom query expression and return a slice of node
func queryFromCssAtom(node *html.Node,cssAtom string)(res []*html.Node){
	attrReg:=regexp.MustCompile(`(\w+)?\([(\w)([^$*]?=)"(\w+)"\])+`)
	switch{
		case cssAtom=="*"{
			res=append(res, Traversal(node,AttrFilter{}))
		}
		case strings.HasPrefix(cssAtom,"#"):{
			res=append(res, Traversal(node,AttrFilter{AttrCondis:[]AttrCondtion{AttrCondtion{"id","=",cssAtom[1:]}}}))
		}
		case strings.HasPrefix(cssAtom,"."):{
			res=append(res, Traversal(node,AttrFilter{AttrCondis:[]AttrCondtion{AttrCondtion{"class","*=",cssAtom[1:]}}}))
		}
		case strings.HasSuffix(cssAtom,":first-child"):{
			return Traversal(node,OrderFilter{Order:1,PositiveSequense:true})
		}
		case strings.HasSuffix(cssAtom,":last-child"):{
			return Traversal(node,OrderFilter{Order:1})
		}
		case strings.HasSuffix(cssAtom,":first-of-type"):{
			return Traversal(node,OrderFilter{Name: strings.TrimSuffix(cssAtom,":first-of-type"),Order:1,PositiveSequense:true})
		}
		case strings.HasSuffix(cssAtom,":last-of-type"):{
			return Traversal(node,OrderFilter{Name: strings.TrimSuffix(cssAtom,":last-of-type"),Order:1})
		}
		case strings.HasSuffix(cssAtom,":only-child"):{
			if node.NextSibling==nil&&node.NextSibling==nil&&&&strings.TrimSuffix(cssAtom,":only-child")==node.Data{
				res=append(res,node)
			}
		}
		case strings.HasSuffix(cssAtom,":empty"):{
			if node.FirstChild==nil&&strings.TrimSuffix(cssAtom,":empty")==node.Data{
				res=append(res,node)
			}
		}
		case 

		case 

		case 

		case 

		case 


		default{
			return Traversal(node,AttrFilter{Name:cssAtom})
		}
	}
}

//Given a node and Filter implementation and return a slice of node.
func Traversal(node *html.Node,filter Filter)(res []*html.Node){
	if filter.Accept(node){
		res=append(res,node)
	}
	if node.FirstChild!=nil{
		res=append(res,Traversal(node.FirstChild,filter.Accept(node.FirstChild)))
	}
	for next:=node.NextSibling;next!=nil;next=next.NextSibling{
		if filter.Accept(next){
			res=append(res,next)
		}
		res=append(res,Traversal(next.FirstChild,filter.Accept(next.FirstChild)))
	}
}

type AttrFilter struct{
	Name string
	AttrCondis []AttrCondtion
}
func (af *AttrFilter)Accept(node *html.Node)bool{
	if af.Name!=""&&node.Data!=af.Name{
		return false
	}
	attrMap=Attribute2Map(node.Attr)
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
func (of *OrderFilter)Accept(node *html.Node)bool{
	o:=1
	if of.PositiveSequense{
		if of.Name!=""{
			for next:=node.PrevSibling;next!=nil;next=next.PrevSibling{
				if next.Type==html.ElementType&&next.Data==of.Name{
					o++
				}
			}
		}else{
			for next:=node.PrevSibling;next!=nil;next=next.PrevSibling{
				if next.Type==html.ElementType{
					o++
				}
			}
		}
	}else{
		if of.Name!=""{
			for next:=node.NextSibling;next!=nil;next=next.NextSibling{
				if next.Type==html.ElementType&&next.Data==of.Name{
					o++
				}
			}
		}else{
			for next:=node.NextSibling;next!=nil;next=next.NextSibling{
				if next.Type==html.ElementType{
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

type Filter Interface{
	func Accept(node *Html.Node)bool
}

func Attribute2Map(attrs []html.Attribute)map[string]string{
	attrMap:=map[string]string
	for _,attr:=range attrs{
		attrMap[attr.Key]=attr.Val
	}
}


//Get all child node of the param node
func ChildElements(node *html.Node)[]*html.Node{
	res:=[]*html.Node{}
	if node.FirstChild!=nil{
		for next:=node.FirstChild;next!=nil;next=next.NextSibling{
			if next.Type=html.ElementType{
				resappend(res,next)
			}
		}
	}
	return res
}
//Get all elements node behind the param node
func BehindElements(node *html.Node)[]*html.Node{
	res:=[]*html.Node{}
	for next:=node.NextSibling;next!=nil;next=next.NextSibling{
		if next.Type=html.ElementType{
			resappend(res,next)
		}
	}
	return res
}