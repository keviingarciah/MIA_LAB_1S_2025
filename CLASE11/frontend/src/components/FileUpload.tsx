import React from "react";

interface FileUploadProps {
  onFileContent: (content: string) => void;
}

const FileUpload: React.FC<FileUploadProps> = ({ onFileContent }) => {
  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        onFileContent(content);
      };
      reader.readAsText(file);
    }
    // Reset the input value after reading the file
    event.target.value = "";
  };

  return (
    <input
      type="file"
      id="file-upload"
      className="hidden"
      onChange={handleFileChange}
      accept=".smia"
    />
  );
};

export default FileUpload;
