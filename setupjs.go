package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func setupProxy_js() {
	var px = []byte(`const fs = require('fs')
const { spawn } = require('child_process')
const { createProxyMiddleware } = require('http-proxy-middleware')

if (fs.existsSync('./unserv/unserv')) {
	const stub = spawn('./unserv', [], { cwd: './unserv/' })
	stub.stdout.on('data', data => {
		console.log(data.toString())
	})
}

const proxy = createProxyMiddleware({ target: 'http://localhost:8080' })

module.exports = function (app) {
	walk('./srv/stub').forEach(url => {
		console.log(url)

		app.use(url, proxy)
	})
}

const walk = (dirPath, filesArr = [], prefix = '/') => {
	fs.readdirSync(dirPath).forEach(file => {
		if (fs.statSync(dirPath + '/' + file).isDirectory()) {
			filesArr = walk(dirPath + '/' + file, filesArr, prefix + file + '/')
		} else {
			const fname = file.split('.')[0]
			if (fname.length > 1) {
				filesArr.push(prefix + fname)
				filesArr.push(prefix + fname + '/')
			}
		}
	})
	return filesArr
}
`)

	if _, err := os.Stat("../src/"); err == nil {
		if _, err := os.Stat("../src/setupProxy.js"); os.IsNotExist(err) {
			ioutil.WriteFile("../src/setupProxy.js", px, 0644)
			fmt.Println("Write file ../src/setupProxy.js")
		} else {
			fmt.Printf("../src/setupProxy.js file exist:\n----------\n%s----------\n\n", px)
		}
	}
}
