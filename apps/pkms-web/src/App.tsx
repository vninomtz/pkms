import dog from "./assets/dog.svg";
import notes from "./assets/notes.json";
import ArrowDown from "./assets/ArrowDown.svg";

function App() {
  return (
    <main className="bg-neutral-50 text-base">
      <header className="m-auto h-screen flex justify-center items-center">
        <div className="flex flex-col items-center gap-2">
          <div className="flex justify-center">
            <img src={dog} />
          </div>
          <p className="text-slate-900/75 w-3/4 sm:w-2/4">
            I write when my mind is full of noise, and I take notes when I find
            something meaningful in books or other sources. Here are some of
            them.
          </p>
          <div>
            <img src={ArrowDown} />
          </div>
        </div>
      </header>
      <section className="px-4">
        <ul className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 auto-rows-[8rem] gap-4">
          {notes.map((n) => (
            <NoteItem key={n.id} note={n} />
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
  let rows = "row-span-1";
  const size = note.html.length;
  if (size <= 100) rows = "row-span-1";
  if (size > 100 && size <= 600) rows = "row-span-2";
  if (size > 600 && size <= 1000) rows = "row-span-3";
  if (size > 1000 && size <= 1500) rows = "row-span-4";
  if (size > 1500) rows = "row-span-5";

  return (
    <li
      className={"flex px-3 py-2 flex-col gap-2 bg-white rounded-2xl " + rows}
    >
      <div className="text-sm text-slate-900 font-bold">{note.title}</div>
      <div
        className="overflow-y-auto prose prose-slate text-wrap break-words"
        dangerouslySetInnerHTML={{ __html: note.html }}
      />
    </li>
  );
};

export default App;
