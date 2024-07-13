import ftReact from '@rabarbra/ft_react';
import { apiClient } from '../api/api_client';
import { API_ENDPOINT } from '../config';

const Imgs = () => {
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
            <div>Imgs</div>
            <form
                enctype="multipart/form-data"
                onSubmit={async (ev) => {
                    ev.preventDefault();
                    const formData = new FormData(ev.target);
                    const resp = await apiClient.post(
                        'img',
                        formData,
                    )
                    if (resp.err)
                        setErr(resp.err)
                    else
                        setImgs(null);
                }}
            >
                <input type='file' name="file" required/>
                <button className="btn" type='submit'>Upload</button>
            </form>
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
}

export default Imgs;