#!node

const fs = require('fs')
let chainTxt = fs.readFileSync('./chainBase.txt/').toString()

function getArgs(s) {
    if (!s) return s
    const inputArg = /([^ ,]+) (func\(([^()]*)\) )?([^ ,]+)/g
    let g = inputArg.exec(s)
    let result = []
    while (g) {
        let [, name, func, funcArgs, funcReturn] = g
        if (name == 'sources') name = `changeTop(${name})`
        if (funcReturn.startsWith('...')) name += '...'
        else if (func) {
            if (funcReturn == "Observable") {
                const inputArg = /(func\(([^()]*)\) )?([^ ,]+)/g
                let args = []
                let g = inputArg.exec(funcArgs)
                while (g) {
                    let [, func1, funcArgs1, funcReturn1] = g
                    args.push([funcReturn1])
                    g = inputArg.exec(funcArgs)
                }
                name = `func(${args.map((x,i)=>'arg'+i+' '+x)}) p.Observable {return ${name}(${args.map((x,i)=>'arg'+i)}).source}`
            }
        } else if (funcReturn == "Observable") name += '.source'
        result.push(name)
        g = inputArg.exec(s)
    }
    return result.join(', ')
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
func (observable *Observable) ${g[1]}(${g[2].replace('Deliver','p.Deliver')}) *Observable {
    return &Observable{p.${g[1]}(${args})(observable.source)}
}
`
        }
        g = regDeliver.exec(content)
    }
})
console.log(chainTxt)
fs.writeFileSync('./chain.go', chainTxt)