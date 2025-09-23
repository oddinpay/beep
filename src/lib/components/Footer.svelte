<script>
  // Optional: pass props for year/company if needed
  const year = new Date().getFullYear();
</script>

<style>
  footer {
    position: relative;
    background: #f5f8fb; 
    color: #1a1f36;
    padding: 60px 40px 30px;
    overflow: hidden;
  }


  footer::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 120px;
    background: #0c2c4b; 
    z-index: -1;
  }

  .footer-content {
    max-width: 1200px;
    margin: auto;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .footer-left {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .footer-left .logo {
    font-weight: bold;
    font-size: 20px;
    color: #2a2d55;
  }

  .footer-left .copyright {
    font-size: 13px;
    color: #555;
  }

  .footer-right {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .footer-right a {
    color: #1a1f36;
    text-decoration: none;
    font-size: 14px;
  }

  .footer-right a:hover {
    text-decoration: underline;
  }
</style>

<footer>
  <div class="footer-content">
    <div class="footer-left">
      <div class="logo">ODDINPAY</div>
      <div class="copyright">© {year} ODDINPAY, LLC.</div>
    </div>

    <!-- Right: Links -->
    <div class="footer-right">
      <a href="#">Visit our full website</a>
      <a href="#">Privacy &amp; terms</a>
    </div>
  </div>
</footer>
