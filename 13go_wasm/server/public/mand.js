function jsmand() {
    const canvas = document.getElementById("cnvs");
    if (canvas.getContext) {
        const ctx = canvas.getContext("2d");
        const startTime = Date.now();
        const w = canvas.width;
        const h = canvas.height;
        const itr = 255;
        var size = 3;
        var arr = [];
        let x, y;
        for (let i=0; i<w; i++) {
            x = (i / w) * size - (size / 2);
            arr[i] = [];
            for (let j=0; j<h; j++) {
                y = (j / h) * size - (size / 2);
                var a = 0;
                var b = 0;
                for (let k=0; k<=itr; k++) {
                    // マンデルブロの計算
                    var _a = a * a - b * b + x;
                    var _b = 2 * a * b + y;
                    a = _a;
                    b = _b;
                    if (a * a + b * b > 4) {
                        break;
                    }
                    arr[i][j] = k;
                }
            }
        }
        const endTime = Date.now();
        for (let i=0; i < w; i++) {
            for (let j=0; j < h; j++) {
                ctx.fillStyle = `hsl(${arr[i][j]}, 100%, 50%)`
                ctx.fillRect(i, j, 1, 1);
            }
        }
        time = endTime - startTime;
        const t = document.getElementById("create-time");
        t.textContent = time + "ミリ秒"
    } else {
        console.log("no context.");
    }
}

// function mand(w, h, itr) {
//     var size = 3;
//     var arr = [];
//     let x, y;
//     for (let i=0; i<w; i++) {
//         x = (i / w) * size - (size / 2);
//         arr[i] = [];
//         for (let j=0; j<h; j++) {
//             y = (j / h) * size - (size / 2);
//             var a = 0;
//             var b = 0;
//             for (let k=0; k<=itr; k++) {
//                 // マンデルブロの計算
//                 var _a = a * a - b * b + x;
//                 var _b = 2 * a * b + y;
//                 a = _a;
//                 b = _b;
//                 if (a * a + b * b > 4) {
//                     break;
//                 }
//                 arr[i][j] = k;
//             }
//         }
//     }
//     return arr;
// }