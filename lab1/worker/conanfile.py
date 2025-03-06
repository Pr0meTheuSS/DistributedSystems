from conan import ConanFile
from conan.tools.cmake import cmake_layout

class StiletConan(ConanFile):
    settings = "os", "compiler", "build_type", "arch"
    generators = ["CMakeDeps", "CMakeToolchain"]

    def requirements(self):
        self.requires("userver/2.7")
        self.requires("sole/1.0.4")
        self.requires("libenvpp/1.5.0")
        self.requires("fmt/10.2.1", override=True)

    def build_requirements(self):
        self.tool_requires("cmake/3.25.1")

    def layout(self):
        cmake_layout(self)
