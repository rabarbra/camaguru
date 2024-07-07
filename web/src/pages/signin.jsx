import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";

const Signin = (props) => {
    return (
        <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
            <form
                className="card-body"
                onSubmit={async (ev)=> {
                    ev.preventDefault();
                    const username = ev.target[0].value
                    const pass = ev.target[1].value
                    const res = await apiClient.authorize({username: username, pass: pass});
                    if (res.ok === "true")
                        props.route("/me");
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
                <span className="label-text">Password</span>
                </label>
                <input type="password" placeholder="password" className="input input-bordered" required />
                <label className="label">
                <a href="#" className="label-text-alt link link-hover">Forgot password?</a>
                </label>
            </div>
            <div className="form-control mt-6">
                <button type="submit" className="btn btn-primary">Login</button>
            </div>
            </form>
        </div>
    );
};

export default Signin;
