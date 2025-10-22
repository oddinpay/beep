import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const data = async () => {
      return (await fetch('http://127.0.0.1:8976/v1/status')).json();
};

  return data();

};
