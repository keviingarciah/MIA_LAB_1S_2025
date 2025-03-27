"use client";
import { useState } from "react";
import InputTerminal from "@/components/InputTerminal";
import OutputTerminal from "@/components/OutputTerminal";
import FileUpload from "@/components/FileUpload";
import { executeCommands } from "@/services/api";

export default function Home() {
  const [input, setInput] = useState("");
  const [output, setOutput] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleClear = () => {
    setInput("");
    setOutput("");
  };

  const handleExecute = async () => {
    if (!input.trim()) {
      setOutput("Por favor ingrese un comando");
      return;
    }

    setIsLoading(true);
    try {
      const result = await executeCommands(input);
      setOutput(result);
    } catch (error) {
      setOutput(error instanceof Error ? error.message : "Error desconocido");
    } finally {
      setIsLoading(false);
    }
  };

  const handleFileContent = (content: string) => {
    setInput(content);
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-8">
      <div className="max-w-5xl mx-auto space-y-6">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
            MIA 1S2025
          </h1>
          <div className="flex gap-4">
            <label
              htmlFor="file-upload"
              className="px-4 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-600 
                       transition-colors duration-200 cursor-pointer shadow-sm
                       flex items-center gap-2"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                />
              </svg>
              Cargar Archivo
            </label>
            <button
              onClick={handleClear}
              className="px-4 py-2 rounded-lg bg-red-500 text-white 
                hover:bg-red-600 transition-colors duration-200 
                shadow-sm flex items-center gap-2"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
              Limpiar
            </button>
            <button
              onClick={handleExecute}
              disabled={isLoading}
              className={`px-4 py-2 rounded-lg bg-green-500 text-white 
                hover:bg-green-600 transition-colors duration-200 
                shadow-sm flex items-center gap-2
                ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
            >
              {isLoading ? (
                <svg className="animate-spin h-5 w-5 mr-2" viewBox="0 0 24 24">
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                    fill="none"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
              ) : (
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              )}
              {isLoading ? "Ejecutando..." : "Ejecutar"}
            </button>
          </div>
        </div>

        <FileUpload onFileContent={handleFileContent} />
        <InputTerminal value={input} onChange={setInput} />
        <OutputTerminal output={output} />
      </div>
    </div>
  );
}
