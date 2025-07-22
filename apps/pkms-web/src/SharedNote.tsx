import { useState } from "react";
import { useLocation, useParams } from "react-router";
import { Decrypt } from "./security";

const useResource = (noteId: string, key: string) => {
  const [content, setContent] = useState("");
  const [error, setError] = useState("");
  const [isOpen, setOpen] = useState(false);

  const getResource = async () => {
    const res = await fetch(
      `${import.meta.env.VITE_API_URL}/resources/${noteId}`,
    );
    setOpen(true);
    if (!res.ok) {
      setError("Nota no existe o ya ha sido eliminada");
      return;
    }
    const body = await res.json();
    const value = await Decrypt(body.result, key);
    setContent(value);
  };

  return {
    content,
    isOpen,
    error,
    open: getResource,
  };
};

const SharedNote = () => {
  const params = useParams();
  const location = useLocation();

  const { content, isOpen, error, open } = useResource(
    params.noteId as string,
    location.hash.slice(1),
  );

  if (!isOpen) {
    return (
      <div className="bg-neutral-50 text-base">
        <div className="m-auto w-2/4 h-screen flex justify-center items-center">
          <button
            className="bg-white text-gray-800 border border-gray-400 rounded-md px-4"
            onClick={open}
          >
            Open note
          </button>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-neutral-50 text-base">
        <div className="m-auto h-screen flex justify-center items-center">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-neutral-50 text-base">
      <div className="m-auto h-screen w-2/4 flex justify-center items-center">
        <TextContent content={content} />
      </div>
    </div>
  );
};

const TextContent = ({ content }: { content: string }) => {
  return (
    <div className="w-full">
      <textarea
        className="w-full border border-gray-200 p-2 rounded-md min-h-5"
        name="notecontent"
        readOnly
        value={content}
      />
    </div>
  );
};

export default SharedNote;
