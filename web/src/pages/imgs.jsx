import ftReact from '@rabarbra/ft_react';
import { apiClient } from '../api/api_client';
import { API_ENDPOINT } from '../config';
import Alert from '../components/alert';

const Imgs = () => {
    const [imgs, setImgs] = ftReact.useState(null);
    const [err, setErr] = ftReact.useState("");
    const [selected, setSelected] = ftReact.useState(-1);
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
                    else {
                        setSelected(-1);
                        setImgs(null);
                    }
                }}
            >
                <input type='file' className='file-input file-input-bordered' name="file" required/>
                <button className="btn" type='submit'>Upload</button>
            </form>
            <div className="grid grid-cols-4 gap-4 mt-3">
            {imgs && imgs.map((item, idx)=>(
                <div className="w-24 h-24" onClick={()=>{
                    setSelected(idx);
                    document.getElementById('my_modal_1').showModal();
                }}>
                    <img src={`${API_ENDPOINT}/${item.link.substring(7)}`} />
                </div>
            ))}
            </div>
            {err && <Alert msg={err}/>}
            <dialog id="my_modal_1" className="modal">
                <div className="modal-box">
                    {selected != -1 && <img src={`${API_ENDPOINT}/${imgs[selected]?.link.substring(7)}`} />}
                    <div className="absolute left-5 right-5 top-1/2 flex -translate-y-1/2 transform justify-between">
                        <a
                            className="btn btn-circle"
                            onClick={()=>{if (selected > 0) setSelected(selected - 1)}}
                            >❮</a>
                        <a
                            className="btn btn-circle"
                            onClick={()=>{if (selected < imgs.length - 1) setSelected(selected + 1)}}
                        >❯</a>
                    </div>
                </div>
                <form method="dialog" className="modal-backdrop">
                    <button>close</button>
                </form>
            </dialog>
        </div>
    )
}

export default Imgs;