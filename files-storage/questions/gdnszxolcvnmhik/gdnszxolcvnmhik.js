import * as fs from 'fs'

class ImportError extends Error { }

const loadModule = async (modulePath) => {
    try {
        return await import(modulePath)
    } catch (e) {
        throw new ImportError(`Unable to import module ${modulePath}`)
    }
}


(async function main() {

    const path = process.argv[2] // submission id
    const param = process.argv[3]

    console.log(path)
    console.log(param)
    // let mod = await loadModule(`../uploaded-code/${path}.js`)

    // try {
    //     const testcases = [process.argv[3]]
    //     for (let i = 0; i < testcases.length; i++) {
    //         let input = testcases[i].input.split("\n").map(v => parseInt(v))
    //         let output = parseInt(testcases[i].output)

    //         const result = mod.sum(input[0], input[1])

    //         const outputPath = `../output-code/${path}`
    //         if (result != output) {
    //             let content = `case no ${i}: failed`
    //             fs.appendFile(`../output-code/${outputPath}.txt`, content, err => {
    //                 if (err) console.log(err)
    //             })
    //         } else {
    //             let content = `case no ${i}: success`
    //             fs.appendFile(`../output-code/${outputPath}.txt`, content, err => {
    //                 if (err) console.log(err)
    //             })
    //         }
    //     }
    // } catch (error) {
    //     console.log("invalid function")
    // } finally {
    //     console.log('finish')
    // }
})()