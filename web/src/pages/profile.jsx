import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";
import { API_ENDPOINT } from "../config";

const Profile = (props) => {
    const [imgs, setImgs] = ftReact.useState(null);
    const [err, setErr] = ftReact.useState("");
    ftReact.useEffect(async () => {
        if (!imgs) {
            const resp = await apiClient.get('img?limit=100');
            if (resp.err)
                setErr(resp.err)
            else
                setImgs(resp);
        }
    }, [imgs, setImgs])
    return (
        <div>
            <div>Profile</div>
            <div className="carousel rounded-box">
            {imgs && imgs.map(item=>(
                <div className="carousel-item">
                    <img src={`${API_ENDPOINT}/${item.link.substring(7)}`} />
                </div>
            ))}
            </div>
            {err && 
                <div role="alert" className="alert alert-error">
                    <span>{err}</span>
                </div>
            }
        </div>
    )
};

export default Profile;