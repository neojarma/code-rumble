import * as fs from "fs"

// import { solveMeFirst } from "./src/2.js"

class ImportError extends Error { }

const loadModule = async (modulePath) => {
    try {
        return await import(modulePath)
    } catch (e) {
        throw new ImportError(`Unable to import module ${modulePath}`)
    }
}

function readDir(dir) {
    return new Promise((resolve, reject) => {
        fs.readdir(dir, (err, files) => {
            if (err) {
                reject(err)
            } else {
                const promises = files.map((v) => {
                    return new Promise((resolve, reject) => {
                        fs.readFile(`${dir}/${v}`, 'utf-8', (err, content) => {
                            if (err) {
                                reject(err)
                            } else {
                                const arr = content.split('\r\n').map((str) => parseInt(str));
                                resolve(arr);
                            }
                        })

                    })
                })

                Promise.all(promises)
                    .then((results) => {
                        resolve(results);
                    })
                    .catch((err) => {
                        reject(err);
                    });
            }
        })
    })
}

async function getCases() {
    var input = await readDir('./questions/idttthqcfgimmqb/input')
    var output = await readDir('./questions/idttthqcfgimmqb/output')

    var cases = []

    for (let i = 0; i < input.length; i++) {

        var obj = [
            `cases no ${i}`,
            ...input[i],
            output[i][0]
        ]

        cases.push(obj)
    }

    return cases

}

(async function main() {
    // var a = await getCases()
    var a = await loadModule('./src/2.js')

    try {
        console.log(a.solveMeFirst(1, 1))
        console.log('asdasda')
    } catch (error) {
        console.log("here ")
    } finally {
        console.log('asd')
    }

    for (let i = 0; i < array.length; i++) {

        { { tobereplaced } }
    }

})()