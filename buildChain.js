const fs = require('fs')
let chainTxt =
    `package rx

import (
	"time"

	p "github.com/langhuihui/gorx/pipe"
)

//Observable emit data
type Observable struct {
	source p.Observable
}

func changeTop(sources []Observable) (pSources []p.Observable) {
	pSources = make([]p.Observable, len(sources))
	for i, source := range sources {
		pSources[i] = source.source
	}
	return
}
//Pipe
func (this *Observable) Pipe(cbs ...p.Deliver) *Observable {
	return &Observable{p.Pipe(this.source, cbs...)}
}
//Subject
func Subject(source Observable, input <-chan interface{}) *Observable {
	return &Observable{p.Subject(source.source, input)}
}
`



function getArgs(s) {
    if (!s) return s
    s = s.replace(/func *\(([^()]+)\)/g, 'func')
    return s.split(', ').map(x => {
        let [name, type, funcType] = x.split(' ')
        if (name == 'sources') name = `changeTop(${name})`
        if (type.startsWith('...')) name += '...'
        else if (type == "Observable") name += '.source'
        else if (type == 'func' && funcType == 'Observable') {
            name = `func(d interface{}) p.Observable {return ${name}(d).source}`
        }
        return name
    }).join(', ')
}
fs.readdirSync('./pipe').forEach(f => {
    let content = fs.readFileSync('./pipe/' + f).toString()
    const regOb = /func ([^()]+)\((.*)\) Observable {/g
    let g = regOb.exec(content)
    while (g) {
        if (g[1] != 'Pipe' && g[1] != 'Subject') {
            let args = getArgs(g[2])
            chainTxt += `
//${g[1]} 
func ${g[1]}(${g[2].replace('Deliver','p.Deliver')}) *Observable {
    return &Observable{p.${g[1]}(${args})}
}
`
        }
        g = regOb.exec(content)
    }
})
fs.readdirSync('./pipe').forEach(f => {
    let content = fs.readFileSync('./pipe/' + f).toString()
    const regDeliver = /func ([^()]+)\((.*)\) Deliver {/g
    let g = regDeliver.exec(content)

    while (g) {
        if (g[1] != 'deliver') {
            let args = getArgs(g[2])
            chainTxt += `
//${g[1]} 
func (this *Observable) ${g[1]}(${g[2].replace('Deliver','p.Deliver')}) *Observable {
    return &Observable{p.${g[1]}(${args})(this.source)}
}
`
        }
        g = regDeliver.exec(content)
    }
})
console.log(chainTxt)
fs.writeFileSync('./chain.go', chainTxt)