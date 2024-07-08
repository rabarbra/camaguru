
import "./public/globals.css"
import ftReact from "@rabarbra/ft_react";
import Router, { Route } from "./router";
import Video from "./pages/video";
import Signin from "./pages/signin";
import Main from "./pages/main";
import Profile from "./pages/profile";
import Header from "./components/header";
import Footer from "./components/footer";
import Signup from "./pages/signup";
import Imgs from "./pages/imgs";

const App = (props) => {
    return (
        <div className="h-screen flex flex-col justify-between">
            <Header/>
            <main className="h-100 flex flex-col items-center">
                <Router>
                    <Route fallback path="/" element={<Main/>}/>
			    	<Route login path="/signin" element={<Signin/>}/>
			    	<Route path="/signup" element={<Signup/>}/>
			    	<Route path="/signup" element={<Video/>}/>
                    <Route path="/video" element={<Video/>}/>
			    	<Route auth path="/me" element={<Profile/>}/>
			    	<Route auth path="/imgs" element={<Imgs/>}/>
                </Router>
            </main>
            <Footer/>
        </div>
    );
};

const root = document.getElementById("app");
ftReact.render(<App/>, root);
