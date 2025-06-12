import dog from "./assets/dog.svg";
import notes from "./assets/notes.json";

function App() {
  return (
    <main className="bg-white text-base">
      <header className="max-w-3xl pt-8 mx-auto px-6 sm:px-12">
        <div className="flex justify-center">
          <img src={dog} />
        </div>
        <div>
          <p className="text-slate-900/75 mt-2">
            I write when my mind is full of noise, and I take notes when I find
            something meaningful in books or other sources. Here are some of
            them.
          </p>
        </div>
      </header>
      <section className="max-w-3xl pt-8 mx-auto px-6 sm:px-12">
        <ul className="flex flex-col gap-4">
          {notes.map((n) => (
            <NoteItem note={n} />
          ))}
        </ul>
      </section>
    </main>
  );
}

interface Note {
  id: string;
  title: string;
  content: string;
  metadata: Metadata;
  html: string;
}
interface Metadata {
  tags: string[];
}
const NoteItem = ({ note }: { note: Note }) => {
  return (
    <li>
      <div className="text-slate-900 font-bold">{note.title}</div>
      <div
        className="prose prose-slate"
        dangerouslySetInnerHTML={{ __html: note.html }}
      ></div>
    </li>
  );
};

export default App;
