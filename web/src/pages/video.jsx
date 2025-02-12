import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";
import Msg from "../components/msg";
import Alert from "../components/alert";

const constraints = {
    audio: false,
    video: { width: 1920, height: 1080 },
};

const Video = (props) => {
    const [err, setErr] = ftReact.useState("");
    const [msg, setMsg] = ftReact.useState("");
    const [sat, setSat] = ftReact.useState(100);
    const [cntr, setCntr] = ftReact.useState(100);
    const [brt, setBrt] = ftReact.useState(100);
    const [hue, setHue] = ftReact.useState(0);
    const [sep, setSep] = ftReact.useState(0);
    const [inv, setInv] = ftReact.useState(0);
    const [gray, setGray] = ftReact.useState(0);
    const [video, setVideo] = ftReact.useState(true);
    const [data, setData] = ftReact.useState("");
    ftReact.useEffect(()=>{
        navigator.mediaDevices
            .getUserMedia(constraints)
            .then((mediaStream) => {
                const video = document.querySelector("video");
                video.srcObject = mediaStream;
            })
            .catch((err) => {
                console.error(`${err.name}: ${err.message}`);
            });
    },[]);
    return (
        <div className="container flex flex-col align-center justify-center gap-2">
            <div
                style={{filter: `
                    hue-rotate(${hue}deg)
                    saturate(${sat}%)
                    sepia(${sep}%)
                    invert(${inv}%)
                    grayscale(${gray}%)
                    contrast(${cntr}%)
                    brightness(${brt}%)
                `}}
            >
                <video
                    className="w-full"
                    style={{display: video ? "block" : "none"}}
                    id="video"
                    onLoadedMetadata={(ev)=>ev.target.play()}
                    ></video>
                <img
                    className="w-full"
                    src={data}
                    style={{display: video ? "none" : "block"}}
                    id="photo"
                ></img>
            </div>
            <canvas
                id="canvas"
                width={constraints.video.width}
                height={constraints.video.height}
                style={{display: 'none'}}
            >
            </canvas>
            <button
                onClick={()=>{
                    if (video) {
                        const canvas = document.getElementById("canvas");
                        // const photo = document.getElementById("photo");
                        const video = document.getElementById("video");
                        const context = canvas.getContext("2d");
                        context.fillStyle = "#AAA";
                        context.filter = "none";
                        context.fillRect(0, 0, canvas.width, canvas.height);
                        context.drawImage(video, 0, 0, constraints.video.width, constraints.video.height);
                        const data = canvas.toDataURL("image/png");
                        setData(data);
                        // photo.setAttribute("src", data);
                        setVideo(false);
                    } else {
                        setVideo(true);
                    }
                }}
                className="btn btn-primary"
            >{video ? "Take photo" : "Camera"}</button>
            <details className="collapse bg-base-200">
                <summary className="collapse-title text-xl font-medium">Show filters</summary>
                <div className="collapse-content">    
                    <label>Saturation: {sat}%</label>
                    <input type="range" min={0} max={200} value={sat} className="range  range-xs" onInput={(ev)=>setSat(ev.target.value)} />
                    <label>Hue: {hue || "0"} deg</label>
                    <input type="range" min={-180} max={180} value={hue} className="range range-xs" onInput={(ev)=>setHue(ev.target.value)} />
                    <label>Contrast: {cntr}%</label>
                    <input type="range" min={0} max={200} value={cntr} className="range range-xs" onInput={(ev)=>setCntr(ev.target.value)} />
                    <label>Brightness: {brt}%</label>
                    <input type="range" min={0} max={200} value={brt} className="range range-xs" onInput={(ev)=>setBrt(ev.target.value)} />
                    <label>Sepia: {sep || "0"}%</label>
                    <input type="range" min={0} max={100} value={sep} className="range range-xs" onInput={(ev)=>setSep(ev.target.value)} />
                    <label>Invert: {inv || "0"}%</label>
                    <input type="range" min={0} max={100} value={inv} className="range range-xs" onInput={(ev)=>setInv(ev.target.value)} />
                    <label>Grayscale: {gray || "0"}%</label>
                    <input type="range" min={0} max={100} value={gray} className="range range-xs" onInput={(ev)=>setGray(ev.target.value)} />
                </div>
            </details>
            <button
                className="btn btn-primary"
                onClick={async () => {
                    const formData = new FormData();
                    const canvas = document.getElementById("canvas");
                    const context = canvas.getContext("2d");
                    context.filter = `
                        hue-rotate(${hue}deg)
                        saturate(${sat}%)
                        sepia(${sep}%)
                        invert(${inv}%)
                        grayscale(${gray}%)
                        contrast(${cntr}%)
                        brightness(${brt}%)
                    `
                    const source = video ? document.getElementById("video") : document.getElementById("photo");
                    context.drawImage(source, 0, 0, constraints.video.width, constraints.video.height);
                    await canvas.toBlob(async (blob)=>{
                        formData.append('file', blob, 'photo.png');
                        const resp = await apiClient.post(
                            'img',
                            formData,
                        );
                        if (resp.err)
                            setErr(resp.err);
                        else if (resp.msg)
                            setMsg(resp.msg);
                    }, "image/png");
                    context.fillStyle = "#AAA";
                    context.filter = "none";
                    context.fillRect(0, 0, canvas.width, canvas.height);
                }}
            >Send</button>
            {msg && <Msg msg={msg}/>}
            {err && <Alert msg={err}/>}
        </div>
    );
};

export default Video;
