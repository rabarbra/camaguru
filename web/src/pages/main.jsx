import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";

const Main = (props) => {
    return (
        <div>
            <h1>Camaguru</h1>
            <form
                enctype="multipart/form-data"
                onSubmit={async (ev) => {
                    ev.preventDefault();
                    const formData = new FormData(ev.target);
                    const resp = await apiClient.post(
                        'img',
                        formData,
                    )
                    console.log(resp);
                }}
            >
                <input type='file' name="file" required/>
                <button className="btn" type='submit'>Upload</button>
            </form>
        </div>
    );
};

export default Main;