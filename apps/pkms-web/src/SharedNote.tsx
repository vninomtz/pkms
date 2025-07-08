import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router";
import { Decrypt } from "./security";

const SharedNote = () => {
  const [content, setContent] = useState("");
  const params = useParams();
  const location = useLocation();

  useEffect(() => {
    fetch(`${import.meta.env.VITE_API_URL}/resources/${params.noteId}`)
      .then((res) => res.json())
      .then((r) => {
        Decrypt(r?.result, location.hash.slice(1)).then((r) => setContent(r));
      });
  }, [params.noteId, location.hash]);
  return (
    <div
      className="overflow-y-auto prose prose-slate text-wrap break-words"
      dangerouslySetInnerHTML={{ __html: content }}
    />
  );
};

export default SharedNote;
