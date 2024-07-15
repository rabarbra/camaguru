import ftReact from '@rabarbra/ft_react';
import { apiClient } from '../api/api_client';
import { API_ENDPOINT } from '../config';
import Alert from '../components/alert';
import ImgInput from '../components/image-imput';

const Imgs = () => {
    const [imgs, setImgs] = ftReact.useState(null);
    const [err, setErr] = ftReact.useState("");
    const [selected, setSelected] = ftReact.useState(-1);
    ftReact.useEffect(async () => {
        if (!imgs && !err) {
            const resp = await apiClient.get('img?limit=100');
            if (resp?.err)
                setErr(resp?.err)
            else if (resp?.length)
                setImgs(resp);
            else
                setErr("No imgs!")
        }
    }, [imgs, setImgs])
    return (
        <div>
            <ImgInput endpoint="img" callback={()=>{
                setImgs(null);
                setSelected(-1);
            }}/>
            <div className="grid grid-cols-4 gap-4 mt-3">
                {imgs && imgs.map((item, idx)=>(
                    <div className="w-24 h-24" onClick={()=>{
                        setSelected(idx);
                        document.getElementById('my_modal_1').showModal();
                    }}>
                        <img src={`${API_ENDPOINT}${item.link}`} />
                    </div>
                ))}
            </div>
            {err && <Alert msg={err}/>}
            <dialog id="my_modal_1" className="modal">
                <div className="modal-box">
                    {selected != -1 && <img src={`${API_ENDPOINT}${imgs[selected].link}`} />}
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