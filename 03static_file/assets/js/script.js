(function() {
    setTimeout(function() {
        const title = document.getElementById("fly");
        title.classList.remove("flyBefore");
        title.classList.add("flyAfter");
    }, 100);

    const img = document.getElementById("fly");
    img.addEventListener("click", function() {
        console.log("fly");
        img.classList.remove("flyAfter");
        img.classList.add("onClickFly");
        setTimeout(function(){
            img.classList.remove("onClickFly");
            img.classList.add("flyBefore");
            setTimeout(function() {
                img.classList.remove("flyBefore");
                img.classList.add("flyAfter");
            }, 10);
        }, 1000);
    });
})();