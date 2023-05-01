// import firebase from "firebase/compat";
//
// export const post = async (url, body) => {
//     const token = await firebase.auth().currentUser.getIdToken();
//     const response = await fetch(url, {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/json',
//             Authorization: `Bearer ${token}`
//         },
//         body: JSON.stringify(body)
//     });
//
//     if (!response.ok) {
//         throw new Error('API Error');
//     }
//
//     return response.json();
// };
