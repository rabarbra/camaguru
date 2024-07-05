
import "./public/globals.css"
import ftReact from "@rabarbra/ft_react";

const constraints = {
    audio: false,
    video: { width: 640, height: 360 },
};

const App = (props) => {
    navigator.mediaDevices
        .getUserMedia(constraints)
        .then((mediaStream) => {
            const video = document.querySelector("video");
            video.srcObject = mediaStream;
            video.onloadedmetadata = () => {
                video.play();
            };
        })
        .catch((err) => {
          // always check for errors at the end.
            console.error(`${err.name}: ${err.message}`);
        });
    return (
        <div className="h-screen bg-slate-700 flex flex-col align-middle justify-center">
            <h1 className="my-auto text-3xl text-center">CAMAGURU</h1>
            <video></video>
        </div>
    );
};

const root = document.getElementById("app");
ftReact.render(<App/>, root);
