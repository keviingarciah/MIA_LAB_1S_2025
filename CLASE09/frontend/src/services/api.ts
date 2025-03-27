const API_URL = "http://localhost:3001";

export const executeCommands = async (command: string): Promise<string> => {
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
