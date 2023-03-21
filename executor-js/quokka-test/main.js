// import fetch from "node-fetch";

// const response = await fetch('http://localhost:8081/question/case?id=gdnszxolcvnmhik&limit=-1')
// const data = await response.json()

// const input = data['data']['testCases'][0]['input']

// console.log(input.split('\\n'))

// import * as fs from 'fs'

// var content = '3'

// fs.appendFile('./log.txt', content, err => {
//     if (err) console.log(err)

//     console.log('ok')
//     console.log(content)
// })

console.log(process.argv[2])
var arr = process.argv[3]
console.log(arr)

console.log(JSON.parse(arr))

// console.log(JSON.parse('[{"id":"cycbjmxhxhpwksl","input":"3\\n2","output":"5"}]'))