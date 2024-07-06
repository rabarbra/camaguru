import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";

const Signin = (props) => {
    return (
        <div className="container flex align-middle justify-center">
            <form
                className="flex flex-col justify-center gap-2"
                onSubmit={(ev)=> {
                    ev.preventDefault();
                    const username = ev.target[0].value
                    const pass = ev.target[1].value
                    apiClient.authorize({username: username, pass: pass});
                }}
            >
                <input type="text" required />
                <input type="password" required />
                <button>LOGIN</button>
            </form>
        </div>
    );
};

export default Signin;