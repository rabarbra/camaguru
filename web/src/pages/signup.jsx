import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";
import Alert from "../components/alert";
import Msg from "../components/msg";

const Signup = (props) => {
    const [err, setErr] = ftReact.useState("");
    const [msg, setMsg] = ftReact.useState("");
    return (
        <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
            {
                msg
                ? <Msg msg={msg}/>
                : <form
                    className="card-body"
                    onSubmit={async (ev)=> {
                        ev.preventDefault();
                        const username = ev.target[0].value
                        const email = ev.target[1].value
                        const pass = ev.target[2].value
                        const res = await apiClient.post("me", {username: username, pass: pass, email: email});
                        if (res.err) {
                            setErr(res.err);
                        } else {
                            setMsg("Email sent. Verify it, please!")
                        }
                    }}
                >
                    <div className="form-control">
                        <label className="label">
                            <span className="label-text">Username</span>
                        </label>
                        <input type="text" placeholder="username" className="input input-bordered" required />
                    </div>
                    <div className="form-control">
                        <label className="label">
                            <span className="label-text">Email</span>
                        </label>
                        <input type="email" placeholder="email" className="input input-bordered" required />
                    </div>
                    <div className="form-control">
                        <label className="label">
                            <span className="label-text">Password</span>
                        </label>
                        <input type="password" placeholder="password" className="input input-bordered" required />
                    </div>
                    <div className="form-control mt-6">
                        <button type="submit" className="btn btn-primary">Create Account</button>
                    </div>
                    {err && <Alert msg={err}/>}
                </form>
            }
        </div>
    );
};

export default Signup;