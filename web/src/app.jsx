
import "./public/globals.css"
import ftReact from "@rabarbra/ft_react";


const App = (props) => {
    return (
        <div className="h-screen bg-slate-700 flex flex-row align-middle justify-center">
            <h1 className="my-auto text-3xl">CAMAGURU</h1>
        </div>
    )
}

const root = document.getElementById("app");
ftReact.render(<App/>, root);
