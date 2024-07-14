import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";
import Alert from "../components/alert";

const Signin = (props) => {
    const [err, setErr] = ftReact.useState("");
    ftReact.useEffect(async () => {
		if (location.search.length) {
			const handleQuery = async () => {
				let queryDict = new URLSearchParams(location.search);
				if (queryDict.has("err")) {
					const msg = queryDict.get("err");
					queryDict.delete("err");
					setErr(msg);
				}
				if (queryDict.has("token")) {
					const res = await apiClient.authorize({
						token: queryDict.get("token"),
					});
					if (res.error)
						setErr(res.error);
					else {
						queryDict.delete("token");
						if (res["ok"] === "true")
							props.route("/");
					}
				}
			}
			await handleQuery();
		}
	}, []);
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
                        props.route("/");
                    else if (res.err)
                        setErr(res.err);
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
                        <a href="/reset" className="label-text-alt link link-hover">Forgot password?</a>
                    </label>
                </div>
                <div className="form-control mt-6">
                    <button type="submit" className="btn btn-primary">Login</button>
                </div>
                {err && <Alert msg={err}/>}
            </form>
        </div>
    );
};

export default Signin;
