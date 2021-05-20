import "regenerator-runtime/runtime";
import 'core-js/stable';
import "./main.css";
import html from "./app.html";

const runtime = require('@wailsapp/runtime');

document.addEventListener('contextmenu', event => event.preventDefault());

function openTab(event, tab) {
    var active = document.getElementsByClassName("show");

    while (active.length) {
        active[0].classList.remove("show");
    }

    document.getElementById(tab).classList.add("show");
    event.currentTarget.classList.add("show");
}

function selectButtonInGroup(event, selector) {
    var active = document.getElementsByClassName("selected " + selector);


    while (active.length) {
        active[0].classList.remove("selected");
    }

    event.currentTarget.classList.add("selected");
}

function setLightMode() {
    document.body.classList.remove("darkmode")

    document.getElementById("mode").innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" enable-background="new 0 0 24 24" height="24px" viewBox="0 0 24 24" width="24px" fill="#000000"><rect fill="none" height="24" width="24"/><path d="M12,3c-4.97,0-9,4.03-9,9s4.03,9,9,9s9-4.03,9-9c0-0.46-0.04-0.92-0.1-1.36c-0.98,1.37-2.58,2.26-4.4,2.26 c-2.98,0-5.4-2.42-5.4-5.4c0-1.81,0.89-3.42,2.26-4.4C12.92,3.04,12.46,3,12,3L12,3z"/></svg>
    `
}

function setDarkMode() {
    document.body.classList.add("darkmode")
    document.getElementById("mode").innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" enable-background="new 0 0 24 24" height="24px" viewBox="0 0 24 24" width="24px" fill="#000000"><rect fill="none" height="24" width="24"/><path d="M12,7c-2.76,0-5,2.24-5,5s2.24,5,5,5s5-2.24,5-5S14.76,7,12,7L12,7z M2,13l2,0c0.55,0,1-0.45,1-1s-0.45-1-1-1l-2,0 c-0.55,0-1,0.45-1,1S1.45,13,2,13z M20,13l2,0c0.55,0,1-0.45,1-1s-0.45-1-1-1l-2,0c-0.55,0-1,0.45-1,1S19.45,13,20,13z M11,2v2 c0,0.55,0.45,1,1,1s1-0.45,1-1V2c0-0.55-0.45-1-1-1S11,1.45,11,2z M11,20v2c0,0.55,0.45,1,1,1s1-0.45,1-1v-2c0-0.55-0.45-1-1-1 C11.45,19,11,19.45,11,20z M5.99,4.58c-0.39-0.39-1.03-0.39-1.41,0c-0.39,0.39-0.39,1.03,0,1.41l1.06,1.06 c0.39,0.39,1.03,0.39,1.41,0s0.39-1.03,0-1.41L5.99,4.58z M18.36,16.95c-0.39-0.39-1.03-0.39-1.41,0c-0.39,0.39-0.39,1.03,0,1.41 l1.06,1.06c0.39,0.39,1.03,0.39,1.41,0c0.39-0.39,0.39-1.03,0-1.41L18.36,16.95z M19.42,5.99c0.39-0.39,0.39-1.03,0-1.41 c-0.39-0.39-1.03-0.39-1.41,0l-1.06,1.06c-0.39,0.39-0.39,1.03,0,1.41s1.03,0.39,1.41,0L19.42,5.99z M7.05,18.36 c0.39-0.39,0.39-1.03,0-1.41c-0.39-0.39-1.03-0.39-1.41,0l-1.06,1.06c-0.39,0.39-0.39,1.03,0,1.41s1.03,0.39,1.41,0L7.05,18.36z"/></svg>
    `
}

function setPaused() {
    document.getElementById("pause").innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" fill="black"><path d="M0 0h24v24H0z" fill="none"/>
        <path d="M8 5v14l11-7z"/>
    </svg>
    `
}

function setResumed() {
    document.getElementById("pause").innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24">
        <path d="M0 0h24v24H0z" fill="none" />
        <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z" />
    </svg>
    `
}

function stopSelect(event) {
    event.preventDefault();
}

function run() {
    let mRate = parseInt(document.getElementById("mutations").value);
    let mAmount = parseFloat(document.getElementById("mutationamount").value);
    let points = parseInt(document.getElementById("points").value);
    let population = parseInt(document.getElementById("population").value);
    let cutoff = parseInt(document.getElementById("cutoff").value);

    let blockSize = parseInt(document.getElementById("blocksize").value);
    let cacheSize = parseInt(document.getElementById("cachesize").value);
    let threads = parseInt(document.getElementById("threads").value);
    let frametime = parseInt(document.getElementById("frametime").value);

    let type = parseInt(document.getElementById("shapecontainer").getElementsByClassName("selected")[0].dataset.shape)

    backend.Runner.Run(type, mRate, mAmount, points, population, cutoff, blockSize, cacheSize, threads, frametime);
}

function stop() {
    backend.Runner.Stop();
}

function updateCanvasSize(ratio) {
    var canvas = document.getElementById("render");
    var area = document.getElementById("renderarea");

    var maxWidth = area.offsetWidth;
    var maxHeight = area.offsetHeight;
    if (maxWidth / maxHeight > ratio) {
        canvas.width = maxHeight * ratio;
        canvas.height = maxHeight;
    } else {
        canvas.width = maxWidth;
        canvas.height = maxWidth / ratio;
    }
}

async function start() {
    var app = document.getElementById('app');
    app.style.width = '100%';
    app.style.height = '100%';

    app.innerHTML = html

    wails.Events.On("image", (path, data) => {
        let name = path.replace(/^.*[\\\/]/, '');
        document.getElementById("select").innerHTML = `
        <img class="preview" src="${data}"/>
        <p class="filename">${name}</p>
        `;
        document.getElementById("select").classList.add("selected");
        document.getElementById("run").removeAttribute("disabled");
    });

    wails.Events.On("stopped", () => {
        document.getElementById("run").innerHTML = "Start";
        document.getElementById("run").onclick = run;
        document.getElementById("pause").setAttribute("disabled", "true");
    })

    wails.Events.On("running", () => {
        document.getElementById("statstab").click();
        document.getElementById("run").innerHTML = "Stop";
        document.getElementById("run").onclick = stop;
        document.getElementById("save").removeAttribute("disabled");
        document.getElementById("pause").removeAttribute("disabled");
    })

    wails.Events.On("paused", () => setPaused())
    wails.Events.On("resumed", () => setResumed())

    wails.Events.On("darkmode", () => setDarkMode())
    wails.Events.On("lightmode", () => setLightMode())

    updateCanvasSize(1);

    window.addEventListener('resize', () => updateCanvasSize(1));

    document.getElementById("logo").setAttribute('draggable', false);


    document.getElementById("configtab").onclick = (event) => openTab(event, 'inputpanel');
    document.getElementById("advancedtab").onclick = (event) => openTab(event, 'advancedpanel');
    document.getElementById("statstab").onclick = (event) => openTab(event, 'statspanel');


    document.getElementById("select").ondragenter = function (event) {
        event.preventDefault();
        event.target.classList.add("over");
    };

    document.getElementById("select").ondragover = function (event) {
        event.preventDefault();
    };

    document.getElementById("select").ondragleave = function (event) {
        event.preventDefault();
        event.target.classList.remove("over");
    };

    document.getElementById("select").ondrop = function (event) {
        event.preventDefault();
        event.target.classList.remove("over");

        if (event.dataTransfer && event.dataTransfer.files) {
            let name = event.dataTransfer.files[0].name;

            if (event.dataTransfer.files[0].type == "image/png" || event.dataTransfer.files[0].type == "image/jpeg") {
                var reader = new FileReader();
                reader.onload = function () {
                    let data = reader.result.replace(/^[^_]*,/, "");
                    backend.Runner.LoadImage(name, data, reader.result.match(/^[^_]*,/)[0]);
                };
                reader.readAsDataURL(event.dataTransfer.files[0]);
            }
        }
    };

    document.getElementById("select").onclick = () => backend.Runner.SelectImage();

    document.getElementById("run").onclick = run;

    document.getElementById("save").onclick = function () {

        if (document.getElementById("png").classList.contains("selected")) {
            var effect = 0;
            if (document.getElementById("gradient").classList.contains("selected")) {
                effect = 1;
            } else if (document.getElementById("split").classList.contains("selected")) {
                effect = 2;
            }
            backend.Runner.SavePNG(parseFloat(document.getElementById("scale").value), effect);
        } else {
            backend.Runner.SaveSVG();
        }
    }

    document.getElementById("pause").onclick = () => backend.Runner.TogglePause();

    document.getElementById("mode").onclick = () => backend.Runner.ToggleMode();

    setResumed();
    setLightMode();

    document.getElementById("png").onclick = function (event) {
        selectButtonInGroup(event, "format");
        document.getElementById("scale").readOnly = false;
        document.getElementById("scale").classList.remove("disabled");
        document.getElementById("scale").removeEventListener("mousedown", stopSelect)

        document.getElementById("effect").classList.remove("disabled");
        document.getElementById("effect").removeAttribute("disabled");
    };

    document.getElementById("svg").onclick = function (event) {
        selectButtonInGroup(event, "format");
        document.getElementById("scale").readOnly = true;
        document.getElementById("scale").classList.add("disabled");
        document.getElementById("scale").addEventListener("mousedown", stopSelect)

        document.getElementById("effect").classList.add("disabled");
        document.getElementById("effect").setAttribute("disabled", "true");
    };

    document.getElementById("none").onclick = (event) => selectButtonInGroup(event, "effect");
    document.getElementById("gradient").onclick = (event) => selectButtonInGroup(event, "effect");
    document.getElementById("split").onclick = (event) => selectButtonInGroup(event, "effect");

    var shapeButtons = document.getElementsByClassName("shapebutton");

    for (var i = 0; i < shapeButtons.length; i++) {
        shapeButtons[i].onclick = function (event) {
            selectButtonInGroup(event, "shapebutton")
        };
    }

    wails.Events.On("renderData", renderData => {
        let width = renderData.Width;
        let height = renderData.Height;

        updateCanvasSize(width / height);
        var canvas = document.getElementById("render");

        let cW = canvas.width;
        let cH = canvas.height;
        var ctx = canvas.getContext("2d", { alpha: false });

        if (window.devicePixelRatio > 1) {
            var canvasWidth = canvas.width;
            var canvasHeight = canvas.height;

            canvas.width = canvasWidth * window.devicePixelRatio;
            canvas.height = canvasHeight * window.devicePixelRatio;

            canvas.style.width = canvasWidth + "px";
            canvas.style.height = canvasHeight + "px";

            ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
        }

        ctx.globalCompositeOperation = "lighter";

        ctx.clearRect(0, 0, cW, cH);

        for (let i in renderData.Polygons) {
            let c = renderData.Colors[i];
            let t = renderData.Polygons[i].Points;
            ctx.fillStyle = `rgb(${Math.round(c.R * 255)}, ${Math.round(c.G * 255)}, ${Math.round(c.B * 255)})`;
            ctx.beginPath();
            ctx.moveTo(Math.round(t[0].X * cW), Math.round(t[0].Y * cH));
            for (var j = 1; j < t.length; j++) {
                ctx.lineTo(Math.round(t[j].X * cW), Math.round(t[j].Y * cH));
            }
            ctx.closePath();
            ctx.fill();
        }
    });

    wails.Events.On("stats", stats => {
        document.getElementById("generation").innerHTML = stats.Generation;
        document.getElementById("fitness").innerHTML = `${Math.round(stats.BestFitness * 1000000000) / 10000000}%`;
        document.getElementById("time").innerHTML = `${Math.round(stats.TimeForGen / 10000) / 100}ms`;
    });
}

runtime.Init(start);