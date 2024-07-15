import ftReact from "@rabarbra/ft_react";
import Alert from "./alert";
import Msg from "./msg";
import { apiClient } from "../api/api_client";

const ImgInput = (props) => {
    const [img, setImg] = ftReact.useState(null);
    const [err, setErr] = ftReact.useState("");
    const [msg, setMsg] = ftReact.useState("");
    return (
        <div className="flex flex-col gap-2">
            {img && <img src={img}/>}
            {err && <Alert msg={err}/>}
            {msg && <Msg msg={msg}/>}
            <form
                className="flex flex-wrap align-middle justify-center justify-items-center gap-1"
                enctype="multipart/form-data"
                onSubmit={async (ev) => {
                    ev.preventDefault();
                    const formData = new FormData(ev.target);
                    const resp = await apiClient.post(
                        props.endpoint,
                        formData,
                    )
                    if (resp.err) {
                        setErr(resp.err);
                        return ;
                    }
                    if (resp.msg) {
                        setMsg(resp.msg);
                    }
                    props.callback && props.callback();
                }}
            >
                <input
                    type='file'
                    className='file-input file-input-bordered grow'
                    accept='.jpg, .jpeg, .png'
                    name="file"
                    onChange={(ev)=>{
                        err && setErr("");
                        msg && setMsg("");
                        if (ev.target.files[0]) {
                            const reader = new FileReader();
                            reader.onload = ev => setImg(ev.target.result);
                            reader.readAsDataURL(ev.target.files[0]);
                        } else {
                            setImg(null);
                        }
                    }}
                    required
                />
                <button className="btn flex-1" type='submit'>Upload</button>
            </form>
        </div>
    );
};

export default ImgInput;