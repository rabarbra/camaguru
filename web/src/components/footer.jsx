import ftReact from "@rabarbra/ft_react";

const Footer = (props) => {
    return (
        <footer className="footer footer-center bg-base-300 text-base-content p-4">
            <aside>
                <p>{`Copyright © ${new Date().getFullYear()} - All right reserved by `}
                    <a href="https://profile.intra.42.fr/users/psimonen">psimonen</a>
                </p>
            </aside>
        </footer>
    );
};

export default Footer;