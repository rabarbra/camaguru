import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";
import { API_ENDPOINT } from "../config";
import Alert from "../components/alert";
import Msg from "../components/msg";
import ImgInput from "../components/image-imput";

const Profile = (props) => {
    const me = JSON.parse(localStorage.getItem("me"));
    const [err, setErr] = ftReact.useState("");
    const [msg, setMsg] = ftReact.useState("");
    return (
        <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
            <ImgInput endpoint="me/avatar"/>
            <form
                className="card-body"
                onSubmit={async (ev)=> {
                    ev.preventDefault();
                    let data = {};
                    const username = ev.target[0].value
                    if (username.length) {
                        data.username = username
                    }
                    const email = ev.target[1].value
                    if (email.length) {
                        data.email = email
                    }
                    const pass = ev.target[2].value
                    if (pass.length) {
                        data.pass = pass
                    }
                    const res = await apiClient.put('me', data);
                    if (res.msg)
                        setMsg(res.msg)
                    else if (res.err)
                        setErr(res.err);
                }}
            >
                <div className="form-control">
                    <label className="label">
                    <span className="label-text">Username</span>
                    </label>
                    <input type="text" placeholder="username" className="input input-bordered" value={me.username} />
                </div>
                <div className="form-control">
                    <label className="label">
                        <span className="label-text">Email</span>
                    </label>
                    <input type="email" placeholder="email" className="input input-bordered" value={me.email} />
                </div>
                <div className="form-control">
                    <label className="label">
                        <span className="label-text">Password</span>
                    </label>
                    <input type="password" placeholder="password" className="input input-bordered" />
                </div>
                <div className="form-control mt-6">
                    <button type="submit" className="btn btn-primary">Update Profile</button>
                </div>
                {err && <Alert msg={err}/>}
                {msg && <Msg msg={msg}/>}
            </form>
        </div>
    )
};

export default Profile;