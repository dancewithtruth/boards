import { FOOTER_HEIGHT } from '@/constants';
import Link from 'next/link';
import { FaGithub } from 'react-icons/fa';

const Footer = () => (
  <footer className="footer items-center p-10 bg-neutral text-neutral-content" style={{ height: FOOTER_HEIGHT }}>
    <div className="items-center grid-flow-col">
      <Link href="/" className="font-bold text-xl">
        Boards
      </Link>
      <p>Copyright © 2023 - All right reserved</p>
    </div>
    <div className="grid-flow-col gap-4 md:place-self-center md:justify-self-end">
      <a href="https://github.com/Wave-95/boards" target="_blank" className="cursor-pointer">
        <FaGithub size={32} />
      </a>
    </div>
  </footer>
);

export default Footer;
