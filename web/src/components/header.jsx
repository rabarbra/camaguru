import ftReact from "@rabarbra/ft_react";
import { apiClient } from "../api/api_client";
import { API_ENDPOINT } from "../config";

const Header = (props) => {
    const me = JSON.parse(localStorage.getItem("me"));
    return (
        <header>
            <div className="navbar bg-base-100">
                <div className="flex-1">
                    <a className="btn btn-ghost text-xl" href="/">Camaguru</a>
                </div>
                <ul className="menu menu-horizontal">
                    <li>
                        <a href="/video">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-camera"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                        </a>
                    </li>
                    <li>
                        <a href="/imgs">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-images"><path d="M18 22H4a2 2 0 0 1-2-2V6"/><path d="m22 13-1.296-1.296a2.41 2.41 0 0 0-3.408 0L11 18"/><circle cx="12" cy="8" r="2"/><rect width="16" height="16" x="6" y="2" rx="2"/></svg>
                        </a>
                    </li>
                </ul>
                {apiClient.authorized() &&
                    <div className="flex-none">
                        <div className="dropdown dropdown-end">
                        <div
                            tabIndex={0}
                            role="button"
                            className="btn btn-ghost btn-circle avatar"
                        >
                            <div className="w-10 rounded-full">
                            <img
                                alt="Profile"
                                src={me.avatar
                                    ? `${API_ENDPOINT}${me.avatar}`
                                    :"https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg"
                                }/>
                            </div>
                        </div>
                        <ul
                            tabIndex={0}
                            className="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 p-2 shadow">
                            <li><a href="/me">Profile</a></li>
                            <li><a
                                href="/signin"
                                onClick={()=>{
                                    apiClient.unauthorize();
                                }}
                            >Logout</a></li>
                        </ul>
                        </div>
                    </div>
                }
            </div>
        </header>
    );
};

export default Header;