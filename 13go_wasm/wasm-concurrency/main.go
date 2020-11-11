package main

import (
	"fmt"
	"strconv"
	"sync"
	"syscall/js"
	"time"
)

func add(this js.Value, i []js.Value) interface{} {
	console := js.Global().Get("console")
	console.Call("log", i[0].Int()+i[1].Int())
	return nil
}

func mult(this js.Value, args []js.Value) interface{} {
	// 指定したIDの要素の値を取得
	value1 := js.Global().Get("document").Call("getElementById", args[0].String()).Get("value").String()
	value2 := js.Global().Get("document").Call("getElementById", args[1].String()).Get("value").String()
	// 取得した値をint型に変換
	int1, _ := strconv.Atoi(value1)
	int2, _ := strconv.Atoi(value2)
	// 乗算
	ans := int1 * int2
	// 答えを文字列に変換
	s := strconv.Itoa(int1) + "*" + strconv.Itoa(int2) + "=" + strconv.Itoa(ans)
	// liタグの作成
	li := js.Global().Get("document").Call("createElement", "li")
	// liタグに値を設定
	li.Set("textContent", s)
	// ulタグにliタグをアペンドチャイルド
	js.Global().Get("document").Call("getElementById", args[2].String()).Call("appendChild", li)
	return nil
}

var Arr [400][400]int

func gomand(this js.Value, args []js.Value) interface{} {
	var wg sync.WaitGroup
	canvas := js.Global().Get("document").Call("getElementById", "cnvs")
	ctx := canvas.Call("getContext", "2d")
	start := time.Now()
	w := 400
	h := 400
	itr := 255
	size := 3
	for i := 0; i < w; i++ {
		x := (float64(i)/float64(w))*float64(size) - (float64(size) / 2)
		for j := 0; j < h; j++ {
			wg.Add(1)
			y := (float64(j)/float64(h))*float64(size) - (float64(size) / 2)
			go mand(x, y, i, j, itr, &wg)
		}
	}
	wg.Wait()
	end := time.Now()
	// fmt.Println(Arr)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			l := 255 - Arr[i][j]
			hsl := "hsl(" + strconv.Itoa(l) + ", 100%, 50%)"
			ctx.Set("fillStyle", hsl)
			ctx.Call("fillRect", i, j, 1, 1)
		}
	}
	processTime := fmt.Sprintf("%vミリ秒\n", (end.Sub(start)).Milliseconds())
	js.Global().Get("document").Call("getElementById", "create-time").Set("textContent", processTime)
	return nil
}

func mand(x float64, y float64, i int, j int, itr int, w *sync.WaitGroup) {
	a := float64(0)
	b := float64(0)
	for k := 0; k <= itr; k++ {
		aTemp := a*a - b*b + x
		bTemp := 2*a*b + y
		a = aTemp
		b = bTemp
		if a*a+b*b > 4 {
			break
		}
		Arr[i][j] = k
	}
	w.Done()
}

func registerCallbacks() {
	js.Global().Set("add", js.FuncOf(add))
	js.Global().Set("mult", js.FuncOf(mult))
	js.Global().Set("gomand", js.FuncOf(gomand))
}
func main() {
	// チャンネルによって永続化
	c := make(chan struct{}, 0)
	println("Go WebAssembly Initialized")
	registerCallbacks()
	<-c
}
