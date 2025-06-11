import dog from "./assets/dog.svg";
function App() {
  return (
    <main className="bg-white text-base">
      <header className="max-w-3xl pt-8 mx-auto px-6 sm:px-12">
        <div className="flex justify-center">
          <img src={dog} />
        </div>
        <div>
          <h1 className="text-lg font-bold text-slate-900">
            Notes from Victor
          </h1>
          <p className="text-slate-900/75 mt-2">
            I write when my mind is full of noise and take notes when I find
            something useful in books or some interesting source, here are some
            of them.
          </p>
        </div>
      </header>
    </main>
  );
}

export default App;
