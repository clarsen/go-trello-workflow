export default ({ children }) => (
  <main>
    {children}
    <style jsx global>{`
      * {
        font-family: Helvetica, serif;
      }
      body {
        margin: 0;
        padding: 25px 50px;
        background-color: #000;
        color: #fff;
      }
      a {
        color: #22bad9;
      }
      p {
        font-size: 14px;
        line-height: 24px;
      }
      article {
        margin: 0 auto;
        max-width: 650px;
      }
      button {
        align-items: center;
        background-color: #22bad9;
        border: 0;
        color: white;
        display: flex;
        padding: 5px 7px;
      }
      button:active {
        background-color: #1b9db7;
        transition: background-color 0.3s;
      }
      button:focus {
        outline: none;
      }
    `}</style>
  </main>
)
