from pages.base_login_page import BaseLoginPage


class LoginPage(BaseLoginPage):
    # TODO: implement the real login flow for belajar_lintas_ronald.
    PATH = "/login"

    def perform_login(self, username: str, password: str) -> None:
        raise NotImplementedError(
            "Implement perform_login for belajar_lintas_ronald (see pages/base_login_page.py)."
        )
