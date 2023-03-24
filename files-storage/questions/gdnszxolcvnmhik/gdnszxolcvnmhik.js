import * as fs from "fs"

class ImportError extends Error { }

const loadModule = async (modulePath) => {
    try {
        return await import(modulePath)
    } catch (e) {
        throw new ImportError(`Unable to import module ${modulePath}`)
    }
}

(async function main() {
    const submissionId = process.argv[2]
    const mod = await loadModule(`../submitted-code/${submissionId}.js`)
    const testCases = JSON.parse(process.argv[3])

    let fileContentResult = ""

    try {

        testCases.forEach((test, i) => {
            const input = test["input"].split("\n").map(v => parseInt(v))
            const output = parseInt(test["output"])
            const result = mod.sum(input[0], input[1])

            if (result === output) {
                fileContentResult += `${test["id"]}=pass=${result}\n`
            } else {
                if (i == testCases.length - 1) {
                    fileContentResult += `${test["id"]}=failed=${result}`
                } else {
                    fileContentResult += `${test["id"]}=failed=${result}\n`
                }
            }
        })

    } catch (error) {
        fileContentResult.result = "all case failed"

    } finally {
        const outPath = `./js-executor/result-code/${submissionId}.txt`
        fs.writeFileSync(outPath, fileContentResult)
        console.log("finish")
    }
})()