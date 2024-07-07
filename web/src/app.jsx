
import "./public/globals.css"
import ftReact from "@rabarbra/ft_react";
import Router, { Route } from "./router";
import Video from "./pages/video";
import Signin from "./pages/signin";
import Main from "./pages/main";
import Profile from "./pages/profile";

const App = (props) => {
    return (
        <div className="h-screen bg-slate-700 flex align-middle justify-center">
            <Router>
                <Route fallback path="/" element={<Main/>}/>
				<Route login path="/signin" element={<Signin/>}/>
				<Route path="/signup" element={<Video/>}/>
                <Route path="/video" element={<Video/>}/>
				<Route auth path="/me" element={<Profile/>}/>
            </Router>
        </div>
    );
};

const root = document.getElementById("app");
ftReact.render(<App/>, root);
