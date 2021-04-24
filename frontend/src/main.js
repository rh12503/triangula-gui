import "regenerator-runtime/runtime";
import 'core-js/stable';
import "./main.css";
import logo from "./assets/logo.svg";

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

    backend.Runner.Run(mRate, mAmount, points, population, cutoff, blockSize, cacheSize, threads, frametime);
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

    app.innerHTML = `
	    <div id="side">
            <div id="title" class="noselect">
                <image id="logo" src="${logo}"/>
                <span id="titlecontent">Triangula</span>
            </div>
            <hr>
            <div class="controls">
                <button class="noselect" type="button" id="select">Drop an image, or click to select</button>
            </div>
            <div id="info">
                <div id="tabs">
                    <button class="tab show noselect" id="configtab">Basic</button>
                    <button class="tab noselect" id="advancedtab">Advanced</button>
                    <button class="tab noselect" id="statstab">Statistics</button>
                </div>

                <div class="panel show noselect" id="inputpanel">


                    <div class="noselect inputs">
                        <table class="formtable">
                            <tr>
                                <td><label class="subhead">Points</label></td>
                                <td><input type="number" id="points" class="input" value="300" min="0" step="5"></td>
                            </tr>

                            <tr>
                                <td><label class="subhead">Mutations</label></td>
                                <td><input type="number" id="mutations" class="input" value="2" min="0" step="1"></td>
                            </tr>
                            <tr>
                                <td><label class="subhead">Variation</label></td>
                                <td><input type="number" id="mutationamount" class="input" value="0.3" min="0" max="1"
                                        step="0.05"></td>
                            </tr>
                            <tr>
                                <td><label class="subhead">Population</label></td>
                                <td><input type="number" id="population" class="input" value="400" min="0" step="10">
                                </td>
                            </tr>
                            <tr>
                                <td><label class="subhead">Cutoff</label></td>
                                <td><input type="number" id="cutoff" class="input" value="5" min="0" step="1"></td>
                            </tr>
                        </table>
                    </div>
                </div>

                <div class="panel noselect" id="advancedpanel">
                    <div class="noselect inputs">
                        <table class="formtable">
                            <tr>
                                <td>
                                    <p class="subhead">Block Size</p>
                                </td>
                                <td><input type="number" id="blocksize" class="input" value="5" min="0" step="1"></td>
                            </tr>
                            <tr>
                                <td>
                                    <p class="subhead">Cache Size</p>
                                </td>
                                <td><input type="number" id="cachesize" class="input" value="22" min="0" step="1"></td>
                            </tr>
                            <tr>
                                <td>
                                    <p class="subhead">Threads</p>
                                </td>
                                <td><input type="number" id="threads" class="input" value="0" min="0" step="1"></td>
                            </tr>
                            <tr>
                                <td>
                                    <p class="subhead">Time per Frame</p>
                                </td>
                                <td><input type="number" id="frametime" class="input" value="250" min="0" step="50">
                                </td>
                            </tr>
                        </table>
                    </div>
                </div>


                <div class="panel noselect" id="statspanel">
                    <div class="noselect inputs">
                        <p class="statsheader">Generation</p>
                        <p id="generation" class="stats selectable">0</p>
                        <p class="statsheader">Fitness</p>
                        <p id="fitness" class="stats selectable">0.0%</p>
                        <p class="statsheader">Time</p>
                        <p id="time" class="stats selectable">0.00ms</p>
                    </div>
                </div>
            </div>

            <div id="sidebottom">
                <button class="noselect" type="button" id="run" disabled>Start</button>
            </div>
        </div>
        <div id="main">
            <div id="topbar">
                <button class="toptab noselect" id="pause" disabled>
                </button>
                <div id="export">
                    <button class="toptab noselect">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24">

                            <path d="M0 0h24v24H0z" fill="none" />
                            <path d=" M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17
                                11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2z" fill="white" />
                        </svg>
                    </button>
                    <div id="exportinput">
                        <p class="subhead noselect vertical">Format</p>
                        <div class="buttongroup">
                            <button id="png" class="noselect buttoningroup selected format" type="button">PNG</button>
                            <button id="svg" class="noselect buttoningroup format" type="button">SVG</button>
                        </div>
                        <p class="subhead noselect vertical">Scale</p>
                        <div class="inputarea">
                            <input type="number" id="scale" class="input" value="1" min="0" step="1">
                        </div>

                        <p class="subhead noselect vertical">Effect</p>
                        <div id="effect" class="buttongroup">
                            <button id="none" class="noselect buttoningroup selected effect" type="button">None</button>
                            <button id="gradient" class="noselect buttoningroup effect" type="button">Gradient</button>
                            <button id="split" class="noselect buttoningroup effect" type="button">Split</button>
                        </div>

                        <button class="noselect" type="button" id="save" disabled>Save</button>
                    </div>
                </div>
            </div>
            <div id="renderarea">
                <canvas id="render"></canvas>
            </div>
        </div>
	`
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


    setResumed();

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


    wails.Events.On("renderData", renderData => {
        let width = renderData.Width;
        let height = renderData.Height;

        updateCanvasSize(width / height);
        var canvas = document.getElementById("render");

        let cW = canvas.width;
        let cH = canvas.height;
        var ctx = canvas.getContext("2d", {alpha: false});

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

        for (let tri of renderData.Data) {
            let c = tri.Color;
            let t = tri.Triangle.Points;
            ctx.fillStyle = `rgb(${Math.round(c.R * 255)}, ${Math.round(c.G * 255)}, ${Math.round(c.B * 255)})`;
            ctx.beginPath();
            ctx.moveTo(Math.round(t[0].X * cW), Math.round(t[0].Y * cH));
            ctx.lineTo(Math.round(t[1].X * cW), Math.round(t[1].Y * cH));
            ctx.lineTo(Math.round(t[2].X * cW), Math.round(t[2].Y * cH));
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