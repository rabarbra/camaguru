import ftReact from "@rabarbra/ft_react";

const constraints = {
    audio: false,
    video: { width: 640, height: 360 },
};

const Video = (props) => {
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
        <video></video>
    )
};

export default Video;