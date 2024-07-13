import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";

const Reset = (props) => {
    const [err, setErr] = ftReact.useState("");
    const [res, setRes] = ftReact.useState("");
    const [token, setToken] = ftReact.useState("");
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
                    setToken(queryDict.get("token"));
                    queryDict.delete("token");
				}
			}
			await handleQuery();
		}
	}, []);
    if (!token.length) {
        return (
            <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
                <form
                    className="card-body"
                    onSubmit={async (ev)=> {
                        ev.preventDefault();
                        const email = ev.target[0].value
                        const res = await apiClient.post('auth/pass', {email: email});
                        if (res.err)
                            setErr(res.err)
                        else if (res.msg)
                            setRes(res.msg);
                    }}
                >
                    {err && 
                        <div role="alert" className="alert alert-error">
                            <span>{err}</span>
                        </div>
                    }
                    {res ? 
                        <div role="alert" className="alert alert-success">
                            <span>{res}</span>
                        </div>
                        : <div>
                            <div className="form-control">
                                <label className="label">
                                    <span className="label-text">Email</span>
                                </label>
                                <input type="email" placeholder="email" className="input input-bordered" required />
                            </div>
                            <div className="form-control mt-6">
                                <button type="submit" className="btn btn-primary">Send me a link to reset password</button>
                            </div>
                        </div>
                    }
                </form>
            </div>
        )
    }
    return (
        <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
            <form
                className="card-body"
                onSubmit={async (ev)=> {
                    ev.preventDefault();
                    const pass = ev.target[0].value
                    const res = await apiClient.post('auth/reset', {pass: pass}, null, {"Authorization": "Bearer " + token});
                    if (res.err)
                        setErr(res.err)
                    else if (res.msg) {
                        setRes(res.msg);
                        props.route("/signin");
                    }
                }}
            >
            {err && 
                <div role="alert" className="alert alert-error">
                    <span>{err}</span>
                </div>
            }
            {res && 
                <div role="alert" className="alert alert-success">
                    <span>{res}</span>
                </div>
            }
            <div className="form-control">
                <label className="label">
                <span className="label-text"> New password</span>
                </label>
                <input type="password" placeholder="password" className="input input-bordered" required />
            </div>
            <div className="form-control mt-6">
                <button type="submit" className="btn btn-primary">Reset Password</button>
            </div>
            </form>
        </div>
    );
}

export default Reset;