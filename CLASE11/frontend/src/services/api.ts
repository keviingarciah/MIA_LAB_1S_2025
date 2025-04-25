const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const executeCommands = async (command: string): Promise<string> => {
  console.log(API_URL);

  try {
    const response = await fetch(`${API_URL}/execute`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ command }),
    });

    if (!response.ok) {
      throw new Error("Error en la respuesta del servidor");
    }

    const data = await response.json();
    return data.output;
  } catch (error) {
    console.error("Error:", error);
    throw new Error("Error al ejecutar los comandos");
  }
};
