import React from "react";

interface OutputTerminalProps {
  output: string;
}

const OutputTerminal = ({ output }: OutputTerminalProps) => {
  return (
    <div className="rounded-lg overflow-hidden shadow-lg border border-gray-200 dark:border-gray-700">
      <div className="bg-gray-100 dark:bg-gray-800 px-4 py-2 border-b border-gray-200 dark:border-gray-700">
        <div className="flex justify-between items-center">
          <div className="flex space-x-2">
            <div className="w-3 h-3 rounded-full bg-red-500"></div>
            <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
            <div className="w-3 h-3 rounded-full bg-green-500"></div>
          </div>
          <span className="text-xs text-gray-500 dark:text-gray-400">
            Output
          </span>
        </div>
      </div>
      <div className="w-full h-48 bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-4 font-mono overflow-auto">
        <pre className="whitespace-pre-wrap">{output}</pre>
      </div>
    </div>
  );
};

export default OutputTerminal;
