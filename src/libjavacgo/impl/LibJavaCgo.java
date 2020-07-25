package libjavacgo.impl;

import org.graalvm.nativeimage.IsolateThread;
import org.graalvm.nativeimage.c.function.CEntryPoint;
import org.graalvm.nativeimage.c.type.CCharPointer;
import org.graalvm.nativeimage.c.type.CTypeConversion;
import com.oracle.svm.core.c.CConst;

public final class LibJavaCgo {
    @CEntryPoint(name = "java_cgo_str")
    public static @CConst CCharPointer javaCgoStr(IsolateThread thread, @CConst CCharPointer s) {
        String str = CTypeConversion.toJavaString(s);
        CTypeConversion.CCharPointerHolder holder = CTypeConversion.toCString(str);
        CCharPointer value = holder.get();
        return value;
    }
}
